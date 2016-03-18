package pushers

import (

"fmt"
"strings"
"net/http"
"io/ioutil"
"encoding/json"
"bytes"
"os"

)


type Resources struct {
	CPU 		int32 			`json:"CPU"`
	MemoryMB 	int32 			`json:"MemoryMB"`
	DiskMB		int32			`json:"DiskMB"`
	IOPS		int32			`json:"IOPS"`
	Networks	interface {} 	`json:"Networks"`
}

type JavaConfig struct {
	Artifact_source 	string 		`json:"artifact_source"`
	Jvm_options 		[]string 	`json:"jvm_options"`
}

type LogConfig struct {
	MaxFiles 		int32 	`json:"MaxFiles"`
	MaxFileSizeMB	int32 	`json:"MaxFileSizeMB"`
}

type RestartPolicy struct {
	Interval 	int64 	`json:"Interval"`
	Attempts 	int32	`json:"Attempts"`
	Delay 		int64 	`json:"Delay"`
	Mode 		string 	`json:"Mode"`
}

type JavaTask struct {
	Name 			string 			`json:"Name"`
	Driver 			string 			`json:"Driver"`
	Config 			JavaConfig		`json:"Config"`
	Constraints 	interface {}	`json:"Constraints"`
	Env 			interface {} 	`json:"Env"`
	Services 		interface {}	`json:"Services"`
	Resources 		Resources 		`json:"Resources"`
	Meta 			interface {} 	`json:"Meta"`
	KillTimeout  	int64			`json:"KillTimeout"`
	LogConfig 		LogConfig		`json:"LogConfig"`
}

type JavaTaskGroup struct {
	Name 			string 			`json:"Name"`
	Count 			int32 			`json:"Count"`
	Constraints		interface {} 	`json:"Constraints"`
	Tasks 			[]JavaTask	 	`json:"Tasks"`
	RestartPolicy 	RestartPolicy 	`json:"RestartPolicy"`
	Meta 			interface {} 	`json:"Meta"`
}

type Constraint struct {
	LTarget 	string 	`json:"LTarget"`
	RTarget 	string 	`json:"RTarget"`
	Operand 	string 	`json:"Operand"`
}

type Update struct {
	Stagger  	int64 	`json:"Stagger"`
	MaxParallel int 	`json:"MaxParallel"`
}

type JobDescription struct {
	Region 				string 			`json:"Region"`
	ID 					string 			`json:"ID"`
	Name 				string  		`json:"Name"`
	Type 				string 			`json:"Type"`
	Priority 			int 			`json:"Priority"`
	AllAtOnce 			bool			`json:"AllAtOnce"`
	Datacenters			[]string 		`json:"Datacenters"`
	Constraints 		[]Constraint 	`json:"Constraints"`
	TaskGroups 			[]JavaTaskGroup `json:"TaskGroups"`
	Update 				Update 			`json:"Update"`
	Periodic 			interface {} 	`json:"Periodic"`
	Meta				interface {} 	`json:"Meta"`
	Status 				string 			`json:"Status"`
	StatusDescription 	string 			`json:"StatusDescription"`
	CreateIndex			int32			`json:"CreateIndex"`
	ModifyIndex			int32			`json:"ModifyIndex"`
}

type Job struct {
	Job 	JobDescription 	 `json:"Job"`	
}

func NomadJobPusher(type_ string){

	baseUrl 	:= "http://160.80.105.5:4646/v1"
	jobAPI 		:= "jobs"
	jobUrl		:= strings.Join([]string{baseUrl, jobAPI}, "/")

	file_name 	:= strings.Join([]string{type_, "job.json"}, "_")
	path 		:= strings.Join([]string{"/home/conet/workspace/goprojects/src/github.com/pierventre/KARAOKE-Pushers", file_name}, "/")

	file, e 	:= ioutil.ReadFile(path)

    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }

    var job Job
    json.Unmarshal(file, &job)

    client := &http.Client{}

    b, err := json.Marshal(&job)
	if err != nil {
		fmt.Println("Encoding Error")
		return
	}

	req, err := http.NewRequest("PUT", jobUrl, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
    if err != nil {
        fmt.Println("HTTP PUT Error")
    }
    defer resp.Body.Close()

    fmt.Printf("Job %s Pushed and Status:%s\n", job.Job.Name, resp.Status)

}