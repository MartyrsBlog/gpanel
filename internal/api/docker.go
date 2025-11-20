package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

// ContainerInfo 容器信息结构体
type ContainerInfo struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Image      string            `json:"image"`
	Status     string            `json:"status"`
	State      string            `json:"state"`
	Ports      []string          `json:"ports"`
	Created    int64             `json:"created"`
	Labels     map[string]string `json:"labels"`
}

// ImageInfo 镜像信息结构体
type ImageInfo struct {
	ID       string            `json:"id"`
	RepoTags []string          `json:"repo_tags"`
	Size     int64             `json:"size"`
	Created  int64             `json:"created"`
	Labels   map[string]string `json:"labels"`
}

// GetContainers 获取容器列表
func GetContainers(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create Docker client",
			"data":    nil,
		})
		return
	}
	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to list containers",
			"data":    nil,
		})
		return
	}

	var containerList []ContainerInfo
	for _, container := range containers {
		var ports []string
		for _, port := range container.Ports {
			if port.PublicPort > 0 {
				ports = append(ports, strconv.Itoa(int(port.PublicPort))+":"+strconv.Itoa(int(port.PrivatePort))+"/"+port.Type)
			} else {
				ports = append(ports, strconv.Itoa(int(port.PrivatePort))+"/"+port.Type)
			}
		}

		containerInfo := ContainerInfo{
			ID:     container.ID[:12],
			Name:   strings.TrimPrefix(container.Names[0], "/"),
			Image:  container.Image,
			Status: container.Status,
			State:  container.State,
			Ports:  ports,
			Created: container.Created,
			Labels: container.Labels,
		}

		containerList = append(containerList, containerInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    containerList,
	})
}

// StartContainer 启动容器
func StartContainer(c *gin.Context) {
	var req struct {
		ContainerID string `json:"container_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create Docker client",
			"data":    nil,
		})
		return
	}
	defer cli.Close()

	err = cli.ContainerStart(context.Background(), req.ContainerID, types.ContainerStartOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to start container: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Container started successfully",
		"data":    nil,
	})
}

// StopContainer 停止容器
func StopContainer(c *gin.Context) {
	var req struct {
		ContainerID string `json:"container_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create Docker client",
			"data":    nil,
		})
		return
	}
	defer cli.Close()

	err = cli.ContainerStop(context.Background(), req.ContainerID, container.StopOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to stop container: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Container stopped successfully",
		"data":    nil,
	})
}

// RestartContainer 重启容器
func RestartContainer(c *gin.Context) {
	var req struct {
		ContainerID string `json:"container_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create Docker client",
			"data":    nil,
		})
		return
	}
	defer cli.Close()

	err = cli.ContainerRestart(context.Background(), req.ContainerID, container.StopOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to restart container: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Container restarted successfully",
		"data":    nil,
	})
}

// RemoveContainer 删除容器
func RemoveContainer(c *gin.Context) {
	var req struct {
		ContainerID string `json:"container_id" binding:"required"`
		Force       bool   `json:"force"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create Docker client",
			"data":    nil,
		})
		return
	}
	defer cli.Close()

	options := types.ContainerRemoveOptions{
		Force: req.Force,
	}

	err = cli.ContainerRemove(context.Background(), req.ContainerID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to remove container: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Container removed successfully",
		"data":    nil,
	})
}

// GetContainerLogs 获取容器日志
func GetContainerLogs(c *gin.Context) {
	containerID := c.Query("container_id")
	tail := c.DefaultQuery("tail", "100")
	
	if containerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Container ID is required",
			"data":    nil,
		})
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create Docker client",
			"data":    nil,
		})
		return
	}
	defer cli.Close()

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	}

	reader, err := cli.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get container logs: " + err.Error(),
			"data":    nil,
		})
		return
	}
	defer reader.Close()

	// 读取日志内容
	buf := make([]byte, 1024)
	var logs []string
	for {
		n, err := reader.Read(buf)
		if err != nil {
			break
		}
		if n > 0 {
			// Docker 日志每行前面有8字节头部信息，需要去掉
			line := string(buf[8:n])
			logs = append(logs, line)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    logs,
	})
}

// GetImages 获取镜像列表
func GetImages(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create Docker client",
			"data":    nil,
		})
		return
	}
	defer cli.Close()

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to list images",
			"data":    nil,
		})
		return
	}

	var imageList []ImageInfo
	for _, img := range images {
		imageInfo := ImageInfo{
			ID:       img.ID[7:19], // 去掉 sha256: 前缀
			RepoTags: img.RepoTags,
			Size:     img.Size,
			Created:  img.Created,
			Labels:   img.Labels,
		}

		imageList = append(imageList, imageInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    imageList,
	})
}

// PullImage 拉取镜像
func PullImage(c *gin.Context) {
	var req struct {
		ImageName string `json:"image_name" binding:"required"`
		Tag       string `json:"tag"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	if req.Tag == "" {
		req.Tag = "latest"
	}

	fullImageName := req.ImageName + ":" + req.Tag

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create Docker client",
			"data":    nil,
		})
		return
	}
	defer cli.Close()

	reader, err := cli.ImagePull(context.Background(), fullImageName, image.PullOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to pull image: " + err.Error(),
			"data":    nil,
		})
		return
	}
	defer reader.Close()

	// 读取拉取进度
	buf := make([]byte, 1024)
	for {
		_, err := reader.Read(buf)
		if err != nil {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Image pulled successfully",
		"data": gin.H{
			"image": fullImageName,
		},
	})
}

// RemoveImage 删除镜像
func RemoveImage(c *gin.Context) {
	var req struct {
		ImageID string `json:"image_id" binding:"required"`
		Force   bool   `json:"force"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create Docker client",
			"data":    nil,
		})
		return
	}
	defer cli.Close()

	options := types.ImageRemoveOptions{
		Force: req.Force,
	}

	_, err = cli.ImageRemove(context.Background(), req.ImageID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to remove image: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Image removed successfully",
		"data":    nil,
	})
}