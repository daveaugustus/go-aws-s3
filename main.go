package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

var (
	s3session    *s3.S3
	templates    *template.Template
	accessKey    string
	accessSecret string
)

const (
	region        = "ap-south-1"
	HTMLTemplates = `templates\*.html`
)

func init() {
	// AWS Credentials initialization
	accessKey, accessSecret = GetAWSCrtedentials()

	// HTML Template parse
	templates = template.Must(template.ParseGlob(HTMLTemplates))

	// AWS session initialization
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, accessSecret, ""),
	})))
}

// The main function cantains default http HandleFunc and the Server
// which runs at PORT 8080
// Any code even if it's routine placed after ListenAndServe line won't
// execute, since the server will be executing and routing to different
// path
func main() {
	// Gorrila mux handler

	router := mux.NewRouter()

	router.HandleFunc("/", BucketListHandler).Methods("GET")
	router.HandleFunc("/create-bucket", CreateBucketGETHandler).Methods("GET")
	router.HandleFunc("/create-bucket", CreateBucketPOSTHandler).Methods("POST")
	router.HandleFunc("/upload", UploadFileGetHandler).Methods("GET")
	router.HandleFunc("/upload", UploadFilePOSTHandler).Methods("POST")
	router.HandleFunc("/:buck_name", ListObjectHandler)

	fmt.Println("Starting the server on port: 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func BucketListHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := ListBuckets()
	if err != nil {
		log.Fatal("Error while fetching buckets!")
	}

	if err := templates.ExecuteTemplate(w, "index.html", resp); err != nil {
		log.Fatal("Couldn't parse html template: index.html")
	}
}

func CreateBucketGETHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "create_bucket.html", nil); err != nil {
		log.Fatal("Couldn't parse html template:")
	}
}

func CreateBucketPOSTHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bucketname := r.PostForm.Get("bucketname")
	fmt.Println("Bucket name: ", bucketname)

	// Create bucket
	resp, err := s3session.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketname),
		ACL:    aws.String(s3.BucketCannedACLPrivate),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(region),
		},
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println("Bucket name already in use!")
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println("Bucket exists and is owned by you!")
			default:
				fmt.Println("default panic!")
				panic(err)
			}
		}
	}
	fmt.Println(resp)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func UploadFileGetHandler(w http.ResponseWriter, r *http.Request) {
	// Initialization of
	buck, err := ListBuckets()
	if err != nil {
		// TODO: Handle error
	}

	// Get method
	// Pass the list of buckets
	// For html list varables where appropriate bucket can be selected
	if err := templates.ExecuteTemplate(w, "upload.html", buck); err != nil {
		log.Fatal("Couldn't parse html template: upload.html")
	}
}

func UploadFilePOSTHandler(w http.ResponseWriter, r *http.Request) {
	// Initialization of
	buck, err := ListBuckets()
	if err != nil {
		// TODO: Handle error
	}

	r.ParseForm()

	file, fileHeader, err := r.FormFile("myfile")
	if err != nil {
		fmt.Printf("Ther is an err in uploading a file: %s", err.Error())
	}
	bucketName := r.PostForm.Get("bucketname")

	fmt.Println("Uploading: ", fileHeader.Filename)
	fmt.Printf("In %s bucket!\n", bucketName)

	// Uploading to the S3 bucket where bucket can be selected using the drpodown
	// in the html template.
	// TODO: Error check needs to be corrected
	response, err := s3session.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileHeader.Filename),
		ACL:    aws.String(s3.BucketCannedACLPrivate),
	})
	if err != nil {
		log.Printf("Error: while uploading to %s bucket", bucketName)
	}
	fmt.Println(response)
	if err := templates.ExecuteTemplate(w, "upload.html", buck); err != nil {
		log.Fatal("Couldn't parse html template: upload.html")
	}
}

func ListObjectHandler(w http.ResponseWriter, r *http.Request) {
	param1 := r.URL.Query().Get("buck_name")
	fmt.Println(param1)
	response, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("saffron-extension"),
	})
	if err != nil {
		// TODO: Handle the error
	}

	if err := templates.ExecuteTemplate(w, "list_objects.html", response); err != nil {
		log.Fatal("Couldn't parse html template: list_objects.html")
	}
}
