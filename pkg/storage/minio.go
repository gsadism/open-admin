package storage

import (
	"context"
	"fmt"
	minio_client "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/url"
	"os"
	"path"
	"time"
)

type Minio struct {
	Client *minio_client.Client
	Server string
	SSL    bool
}

func NewMinioClient(Endpoint string, AccessKeyID string, SecretAccessKey string, Token string, UseSSL bool) (*Minio, error) {
	client, err := minio_client.New(Endpoint, &minio_client.Options{
		Creds:  credentials.NewStaticV4(AccessKeyID, SecretAccessKey, Token),
		Secure: UseSSL,
	})

	if err != nil {
		return nil, err
	}
	return &Minio{
		Client: client,
		Server: Endpoint,
		SSL:    UseSSL,
	}, nil
}

// ContentType : 获取文件传输类型
func (this *Minio) ContentType(suf string) string {
	return ContentType[suf]
}

// SetBucketPolicy : 设置桶的访问权限
// @params ctx: 上下文
// @params BucketName: 桶名称
// @params policy: 策略(public、private)
// @return error:
func (this *Minio) SetBucketPolicy(ctx context.Context, BucketName string, policy string) error {
	var str string
	if policy == "public" {
		// 设置为public权限
		str = "{\"Version\":\"2012-10-17\"," +
			"\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":" +
			"{\"AWS\":[\"*\"]},\"Action\":[\"s3:ListBucket\",\"s3:ListBucketMultipartUploads\"," +
			"\"s3:GetBucketLocation\"],\"Resource\":[\"arn:aws:s3:::" + BucketName +
			"\"]},{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:PutObject\",\"s3:AbortMultipartUpload\",\"s3:DeleteObject\",\"s3:GetObject\",\"s3:ListMultipartUploadParts\"],\"Resource\":[\"arn:aws:s3:::" +
			BucketName +
			"/*\"]}]}"
	} else {
		// 设置为私有权限
		str = "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
	}

	err := this.Client.SetBucketPolicy(ctx, BucketName, str)
	if err != nil {
		return err
	}
	return nil
}

// UploadByName ： 上传文件(文件路径)
// @params ctx: 上下文
// @params BucketName: 桶名称
// @params policy: 策略(public、private)
// @params Location: 位置
// @return info: UploadInfo contains information about the newly uploaded or copied object.
// @return error:
func (this *Minio) UploadByName(ctx context.Context, BucketName string, FileName string, FilePath string) (minio_client.UploadInfo, error) {
	// Upload the zip file with FPutObject
	suf := path.Ext(FileName) // 文件后缀
	if info, err := this.Client.FPutObject(ctx, BucketName, FileName, FilePath, minio_client.PutObjectOptions{
		ContentType: this.ContentType(suf),
	}); err != nil {
		return minio_client.UploadInfo{}, err
	} else {
		return info, nil
	}
}

func (this *Minio) UploadByByte(ctx context.Context, File io.Reader, ObjectSize int64, BucketName string, ObjectName string) (minio_client.UploadInfo, error) {

	suf := path.Ext(ObjectName) // 文件后缀
	if info, err := this.Client.PutObject(ctx, BucketName, ObjectName, File, ObjectSize, minio_client.PutObjectOptions{ContentType: this.ContentType(suf)}); err != nil {
		return minio_client.UploadInfo{}, err
	} else {
		return info, nil
	}

}

// Download : 下载文件到本地
func (this *Minio) Download(ctx context.Context, BucketName string, FileName string, path string) error {
	object, err := this.Client.GetObject(ctx, BucketName, FileName, minio_client.GetObjectOptions{})
	if err != nil {
		return err
	}
	localFile, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(localFile, object)
	if err != nil {
		return err
	}
	return nil
}

// DeleteBucket2File : 删除指定bucket下的文件
func (this *Minio) DeleteBucket2File(ctx context.Context, BucketName string, FileName string) error {
	err := this.Client.RemoveObject(ctx, BucketName, FileName, minio_client.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// DeleteBucket : 删除指定bucket
func (this *Minio) DeleteBucket(ctx context.Context, BucketName string) error {
	err := this.Client.RemoveBucket(ctx, BucketName)
	if err != nil {
		return err
	}
	return nil
}

// ListBuckets : 获取所有bucket
func (this *Minio) ListBuckets(ctx context.Context) ([]minio_client.BucketInfo, error) {
	buckets, err := this.Client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	return buckets, nil
}

// ListBucketObjects : 获取指定bucket下的所有文件
func (this *Minio) ListBucketObjects(ctx context.Context, BucketName string) []minio_client.ObjectInfo {
	var data []minio_client.ObjectInfo
	data = make([]minio_client.ObjectInfo, 1)

	for item := range this.Client.ListObjects(ctx, BucketName, minio_client.ListObjectsOptions{}) {
		data = append(data, item)
	}

	return data
}

// PresignedGetObject : TODO 获取签名（临时文件访问路径）
/**
* @author: 羽
* @date: 4:27 2023/4/25
* @param: BucketName 桶名称
* @param: ObjectName 对象名称
* @param: Ex 过期时间,至少30s,小于30s,则默认30s
* @return:
* @description:
**/
func (this *Minio) PresignedGetObject(ctx context.Context, BucketName string, ObjectName string, Ex time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%v\"", ObjectName))
	sign, err := this.Client.PresignedGetObject(ctx, BucketName, ObjectName, Ex, reqParams)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

// StatObject ： todo 判断某个文件是否存在
func (this *Minio) StatObject(ctx context.Context, bucketName string, ObjectName string) bool {
	if _, err := this.Client.StatObject(ctx, bucketName, ObjectName, minio_client.StatObjectOptions{}); err != nil {
		if minio_client.ToErrorResponse(err).Code == "NoSuchKey" {
			return false
		} else {
			return false
		}
	}
	return true
}

// BucketExists : todo 判断某个bucket是否存在
func (this *Minio) BucketExists(ctx context.Context, BucketName string) bool {
	if ok, err := this.Client.BucketExists(ctx, BucketName); err != nil || !ok {
		return false
	}
	return true
}

// MakeBucket : todo 创建一个bucket
func (this *Minio) MakeBucket(ctx context.Context, BucketName string, Location string) error {
	if err := this.Client.MakeBucket(ctx, BucketName, minio_client.MakeBucketOptions{
		Region: Location,
	}); err != nil {
		return err
	}
	return nil
}
