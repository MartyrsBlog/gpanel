package api

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gpanel/internal/config"

	"github.com/gin-gonic/gin"
)

// FileInfo 文件信息结构体
type FileInfo struct {
	Name         string  `json:"name"`
	Path         string  `json:"path"`
	Size         int64   `json:"size"`
	IsDir        bool    `json:"is_dir"`
	ModTime      string  `json:"mod_time"`
	Permissions  string  `json:"permissions"`
	Owner        string  `json:"owner"`
	Group        string  `json:"group"`
}

// GetFileList 获取文件列表
func GetFileList(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	// 安全检查，防止路径遍历攻击
	if !isValidPath(path) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid path",
			"data":    nil,
		})
		return
	}

	fileInfos, err := scanDirectory(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    fileInfos,
	})
}

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	path := c.PostForm("path")
	if path == "" {
		path = "/"
	}

	// 安全检查
	if !isValidPath(path) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid path",
			"data":    nil,
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "No file uploaded",
			"data":    nil,
		})
		return
	}
	defer file.Close()

	// 确保目标目录存在
	if err := os.MkdirAll(path, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create directory",
			"data":    nil,
		})
		return
	}

	// 创建目标文件
	filePath := filepath.Join(path, header.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create file",
			"data":    nil,
		})
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to save file",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "File uploaded successfully",
		"data": gin.H{
			"filename": header.Filename,
			"size":     header.Size,
			"path":     filePath,
		},
	})
}

// DownloadFile 下载文件
func DownloadFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Path is required",
			"data":    nil,
		})
		return
	}

	// 安全检查
	if !isValidPath(path) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid path",
			"data":    nil,
		})
		return
	}

	// 检查文件是否存在
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "File not found",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Failed to access file",
				"data":    nil,
			})
		}
		return
	}

	// 如果是目录，创建ZIP文件
	if fileInfo.IsDir() {
		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", filepath.Base(path)))
		
		if err := createZipArchive(c.Writer, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Failed to create archive",
				"data":    nil,
			})
			return
		}
	} else {
		// 单个文件
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))
		c.File(path)
	}
}

// ExtractFile 解压文件
func ExtractFile(c *gin.Context) {
	var req struct {
		FilePath string `json:"file_path" binding:"required"`
		ToPath   string `json:"to_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	// 安全检查
	if !isValidPath(req.FilePath) || !isValidPath(req.ToPath) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid path",
			"data":    nil,
		})
		return
	}

	// 确保目标目录存在
	if err := os.MkdirAll(req.ToPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create directory",
			"data":    nil,
		})
		return
	}

	// 根据文件扩展名选择解压方式
	ext := strings.ToLower(filepath.Ext(req.FilePath))
	switch ext {
	case ".zip":
		err := extractZip(req.FilePath, req.ToPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Failed to extract ZIP file",
				"data":    nil,
			})
			return
		}
	case ".gz", ".tgz":
		err := extractGzip(req.FilePath, req.ToPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Failed to extract GZIP file",
				"data":    nil,
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Unsupported file format",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "File extracted successfully",
		"data":    nil,
	})
}

// ChmodFile 修改文件权限
func ChmodFile(c *gin.Context) {
	var req struct {
		Path  string `json:"path" binding:"required"`
		Mode  string `json:"mode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	// 安全检查
	if !isValidPath(req.Path) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid path",
			"data":    nil,
		})
		return
	}

	// 解析权限模式
	mode, err := strconv.ParseUint(req.Mode, 8, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid permission mode",
			"data":    nil,
		})
		return
	}

	// 修改文件权限
	err = os.Chmod(req.Path, os.FileMode(mode))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to change file permissions",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "File permissions changed successfully",
		"data":    nil,
	})
}

// RenameFile 重命名文件
func RenameFile(c *gin.Context) {
	var req struct {
		OldPath string `json:"old_path" binding:"required"`
		NewPath string `json:"new_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	// 安全检查
	if !isValidPath(req.OldPath) || !isValidPath(req.NewPath) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid path",
			"data":    nil,
		})
		return
	}

	// 重命名文件
	err := os.Rename(req.OldPath, req.NewPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to rename file",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "File renamed successfully",
		"data":    nil,
	})
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	// 安全检查
	if !isValidPath(req.Path) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid path",
			"data":    nil,
		})
		return
	}

	// 删除文件或目录
	err := os.RemoveAll(req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete file",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "File deleted successfully",
		"data":    nil,
	})
}

// scanDirectory 扫描目录
func scanDirectory(path string) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileInfos []FileInfo
	for _, entry := range entries {
		info, _ := entry.Info()
		fullPath := filepath.Join(path, entry.Name())

		fileInfo := FileInfo{
			Name:        entry.Name(),
			Path:        fullPath,
			Size:        info.Size(),
			IsDir:       entry.IsDir(),
			ModTime:     info.ModTime().Format("2006-01-02 15:04:05"),
			Permissions: info.Mode().String(),
		}

		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}

// isValidPath 检查路径是否安全
func isValidPath(path string) bool {
	// 防止路径遍历攻击
	return !strings.Contains(path, "..") && !strings.HasPrefix(path, "/etc") && 
	       !strings.HasPrefix(path, "/root") && !strings.HasPrefix(path, "/sys") &&
	       !strings.HasPrefix(path, "/proc")
}

// createZipArchive 创建ZIP压缩包
func createZipArchive(w io.Writer, source string) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/"
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// extractZip 解压ZIP文件
func extractZip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)
		
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.FileInfo().Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

// extractGzip 解压GZIP文件
func extractGzip(src, dest string) error {
	reader, err := os.Open(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	filename := filepath.Base(src)
	if strings.HasSuffix(filename, ".gz") {
		filename = strings.TrimSuffix(filename, ".gz")
	}

	targetPath := filepath.Join(dest, filename)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, gzReader)
	return err
}