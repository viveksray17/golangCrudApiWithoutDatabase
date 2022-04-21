package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Model for Course - file
type Course struct {
	CourseId     string  `json:"courseid"`
	CourseName   string  `json:"coursename"`
	CoursePrice  int     `json:"price"`
	CourseAuthor *Author `json:"author"`
}

type Author struct {
	FullName string `json:"fullname"`
	Wesite   string `json:"website"`
}

// Fake Database
var courses []Course

// middleware / helper functions - file
func (course Course) IsEmpty() bool {
	return course.CourseName == ""
}

// controllers - file
// server home route
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to API by learncodeonline</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get All Courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one course")
	w.Header().Set("Content-Type", "application/json")
	// grab id from request
	params := mux.Vars(r)

	// loop through courses, find matching id and return the response
	for _, course := range courses {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode(fmt.Sprintf("No course found with given id %v", params["id"]))
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create One Course")
	w.Header().Set("Content-Type", "application/json")

	// whatif: body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}
	// whatif: only {} is provided
	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("No data inside the JSON")
		return
	}

	// whatif: course already exists
	for _, coursefromloop := range courses {
		if coursefromloop.CourseName == course.CourseName {
			json.NewEncoder(w).Encode(fmt.Sprintf("Course %v already exists", course.CourseName))
			return
		}
	}
	// generate unique id, string
	// append new course into courses
	rand.Seed(time.Now().UnixNano())
	course.CourseId = strconv.Itoa(rand.Intn(100)) // random number from 0 to 100
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update One Course")
	w.Header().Set("Content-Type", "application/json")

	// first - grab id from request
	params := mux.Vars(r)

	// loop through the database, find the id, delete that course, create a new course with the same id
	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseId = params["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
		//TODO: send a response when id is not found
	}
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete One Course")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			return
		}
	}
}

func main() {
	fmt.Println("API - learncodeonline.in")
	r := mux.NewRouter()
	courses = append(courses, Course{CourseId: "2", CourseName: "ReactJS", CoursePrice: 299, CourseAuthor: &Author{FullName: "Hitesh Choudhary", Wesite: "learncodeonline.in"}})
	courses = append(courses, Course{CourseId: "4", CourseName: "MernStack", CoursePrice: 199, CourseAuthor: &Author{FullName: "Hitesh Choudhary", Wesite: "go.dev"}})

	// routing
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/course/{id}", deleteOneCourse).Methods("DELETE")
	// listen to a port
	log.Fatal(http.ListenAndServe(":4000", r))
}
