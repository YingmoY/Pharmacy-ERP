package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/file/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UploadFileRequest 文件上传请求。
type UploadFileRequest struct {
	File         multipart.File
	FileHeader   *multipart.FileHeader
	BusinessType string
	BusinessID   string
	UploaderID   int64
	Remark       *string
}

// FileInfoDTO 文件元数据传输对象。
type FileInfoDTO struct {
	FileID       string  `json:"file_id"`
	OriginalName string  `json:"original_name"`
	FileSize     int64   `json:"file_size"`
	ContentType  string  `json:"mime_type"`
	BusinessType string  `json:"business_type,omitempty"`
	BusinessID   string  `json:"business_id,omitempty"`
	UploaderID   int64   `json:"uploader_id"`
	Status       int8    `json:"status"`
	CreatedAt    string  `json:"created_at"`
}

// FileService 定义文件业务逻辑接口。
type FileService interface {
	// UploadFile 接收 multipart 文件，存储到磁盘并写入元数据。
	UploadFile(ctx context.Context, req UploadFileRequest) (*FileInfoDTO, error)
	// GetFileInfo 根据 file_id 获取文件元数据。
	GetFileInfo(ctx context.Context, fileID string) (*FileInfoDTO, error)
	// GetFilePath 根据 file_id 获取文件在磁盘上的绝对路径，用于流式下载。
	GetFilePath(ctx context.Context, fileID string) (string, string, error)
}

type fileService struct {
	db      *gorm.DB
	log     *zap.Logger
	baseDir string // 文件存储根目录，默认 ./uploads
}

// NewFileService 创建文件服务实例。
// baseDir 为空时默认使用 "./uploads"。
func NewFileService(db *gorm.DB, log *zap.Logger, baseDir string) FileService {
	if baseDir == "" {
		baseDir = "./uploads"
	}
	return &fileService{
		db:      db,
		log:     log,
		baseDir: baseDir,
	}
}

// toDTO 将 GORM 模型转换为 DTO。
func toDTO(m *model.FileInfo) *FileInfoDTO {
	return &FileInfoDTO{
		FileID:       m.FileID,
		OriginalName: m.OriginalName,
		FileSize:     m.FileSize,
		ContentType:  m.ContentType,
		BusinessType: m.BusinessType,
		BusinessID:   m.BusinessID,
		UploaderID:   m.UploaderID,
		Status:       m.Status,
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// UploadFile 上传文件：写磁盘 + 写数据库元数据。
func (s *fileService) UploadFile(ctx context.Context, req UploadFileRequest) (*FileInfoDTO, error) {
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")

	// 构造存储目录：{baseDir}/{year}/{month}/
	dir := filepath.Join(s.baseDir, year, month)
	if err := os.MkdirAll(dir, 0755); err != nil {
		s.log.Error("创建文件存储目录失败", zap.String("dir", dir), zap.Error(err))
		return nil, ecode.ErrSystem
	}

	// 生成唯一文件 ID 和存储文件名
	fileID := uuid.New().String()
	ext := filepath.Ext(req.FileHeader.Filename)
	storageName := fileID + ext
	storagePath := filepath.Join(dir, storageName)

	// 写文件到磁盘并计算 MD5
	dst, err := os.Create(storagePath)
	if err != nil {
		s.log.Error("创建目标文件失败", zap.String("path", storagePath), zap.Error(err))
		return nil, ecode.ErrSystem
	}
	defer dst.Close()

	hasher := md5.New()
	writer := io.MultiWriter(dst, hasher)
	written, err := io.Copy(writer, req.File)
	if err != nil {
		s.log.Error("写入文件失败", zap.String("path", storagePath), zap.Error(err))
		_ = os.Remove(storagePath) // 写入失败则清理
		return nil, ecode.ErrSystem
	}

	fileHash := fmt.Sprintf("%x", hasher.Sum(nil))
	contentType := req.FileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 写入数据库元数据
	info := &model.FileInfo{
		FileID:       fileID,
		OriginalName: req.FileHeader.Filename,
		StorageName:  storageName,
		StoragePath:  storagePath,
		ContentType:  contentType,
		FileSize:     written,
		FileHash:     fileHash,
		BusinessType: req.BusinessType,
		BusinessID:   req.BusinessID,
		UploaderID:   req.UploaderID,
		Status:       model.FileStatusActive,
		Remark:       req.Remark,
	}

	if err := s.db.WithContext(ctx).Create(info).Error; err != nil {
		s.log.Error("保存文件元数据失败", zap.Error(err))
		_ = os.Remove(storagePath) // 数据库写入失败则清理磁盘文件
		return nil, err
	}

	s.log.Info("文件上传成功",
		zap.String("file_id", fileID),
		zap.String("original_name", req.FileHeader.Filename),
		zap.Int64("size", written),
	)

	return toDTO(info), nil
}

// GetFileInfo 根据 file_id 获取文件元数据。
func (s *fileService) GetFileInfo(ctx context.Context, fileID string) (*FileInfoDTO, error) {
	var info model.FileInfo
	if err := s.db.WithContext(ctx).Where("file_id = ?", fileID).First(&info).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.New(40403, "file not found")
		}
		return nil, err
	}
	return toDTO(&info), nil
}

// GetFilePath 根据 file_id 获取文件磁盘路径与原始文件名，用于下载。
func (s *fileService) GetFilePath(ctx context.Context, fileID string) (string, string, error) {
	var info model.FileInfo
	if err := s.db.WithContext(ctx).Where("file_id = ?", fileID).First(&info).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ecode.New(40403, "file not found")
		}
		return "", "", err
	}

	// 检查文件是否存在于磁盘
	if _, err := os.Stat(info.StoragePath); os.IsNotExist(err) {
		s.log.Error("文件磁盘记录丢失", zap.String("file_id", fileID), zap.String("path", info.StoragePath))
		return "", "", ecode.New(60001, "file storage missing")
	}

	return info.StoragePath, info.OriginalName, nil
}
