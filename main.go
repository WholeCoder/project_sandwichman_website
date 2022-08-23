package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func indexHandler(writer http.ResponseWriter, request *http.Request, nameOfFileWithPath string) {
	writer.Header().Set("Content-Type", "application/octet-stream")

	html, err := template.ParseFiles(nameOfFileWithPath)
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}

// Compile templates on start of the application
var templates = template.Must(template.ParseFiles("public/upload.html"))

// Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) string {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return ""
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create("./uploads/" + handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return "./uploads/" + handler.Filename
}

func decompressHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		fileWithPath := uploadFile(w, r)
		fmt.Println("----------------------------------->", getOnlyFileNameWithExtension(fileWithPath))
		decompress_main(fileWithPath, getOnlyFileNameWithExtension(fileWithPath))
		redirectToFreshman(w, r, getOnlyFileNameWithExtension(fileWithPath))
	}
}

func getOnlyFileNameWithExtension(fileWithPath string) string {
	fmt.Println("fileWithpath ==", fileWithPath)
	splitAtDots := strings.Split(fileWithPath, ".")
	fmt.Println("Split by . ------------->", splitAtDots)
	correctExtension := splitAtDots[1] + "." + splitAtDots[2]

	fmt.Println("correctExtension = ", correctExtension)

	splitAtSlashes := strings.Split(correctExtension, "/")
	fmt.Println("splitAtSlashes[2] == ", splitAtSlashes[2])
	return splitAtSlashes[2]
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		fileWithPath := uploadFile(w, r)
		compress_main(fileWithPath, fileWithPath+".comp")
		redirectToFreshman(w, r, fileWithPath+".comp")
	}
}

func redirectToFreshman(w http.ResponseWriter, r *http.Request, fileWithPath string) {

	w.Header().Set("Content-Type", "application/octet-stream")

	v := strings.Split(fileWithPath, "/")
	fmt.Println("------------->", "/downloads/"+v[0])
	http.Redirect(w, r, "/downloads/"+v[0], http.StatusSeeOther)
}

func main() {
	// Upload route
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/decompress", decompressHandler)

	//	http.HandleFunc("/download", redirectToFreshman)
	fsHandler := http.FileServer(http.Dir("./uploads"))
	http.Handle("/downloads/", http.StripPrefix("/downloads", fsHandler))

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
