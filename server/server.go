package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"

	"../env"
	"../obj"
)

func main() {
	go server()

	var finish string
	fmt.Scanln(&finish)
}

func server() {
	server := new(Server)
	server.Students = make(map[string]map[string]float64)
	server.Subjects = make(map[string]map[string]float64)

	rpc.Register(server)

	listener, err := net.Listen(env.ConnectionType, env.ConnectionPort)
	if err != nil {
		fmt.Println("LISTENER ERROR:", err)
		return
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("CONNECTION ERROR:", err)
			continue
		}

		go rpc.ServeConn(connection)
	}
}

type Server struct {
	Students map[string]map[string]float64
	Subjects map[string]map[string]float64
}

func (svr *Server) SetGradeStudent(data obj.ContainerFull, reply *string) error {
	_, err := svr.Students[data.Student][data.Subject]
	if err {
		return errors.New("the grade ​has already been registered")
	}
	_, err = svr.Subjects[data.Subject][data.Student]
	if err {
		return errors.New("the grade ​has already been registered")
	}

	_, ok := svr.Students[data.Student]
	if !ok {
		svr.Students[data.Student] = map[string]float64{data.Subject: data.Grade}
	} else {
		svr.Students[data.Student][data.Subject] = data.Grade
	}

	_, ok = svr.Subjects[data.Subject]
	if !ok {
		svr.Subjects[data.Subject] = map[string]float64{data.Student: data.Grade}
	} else {
		svr.Subjects[data.Subject][data.Student] = data.Grade
	}

	*reply = "the student has been registered"
	return nil
}

func (svr *Server) GetAverageStudent(data obj.ContainerStudent, reply *string) error {
	_, err := svr.Students[data.Student]
	if !err {
		return errors.New("the student has not been registered")
	}

	var average float64 = 0
	for _, grade := range svr.Students[data.Student] {
		average += grade
	}
	average /= float64(len(svr.Students[data.Student]))

	*reply = fmt.Sprint(average)
	return nil
}

func (svr *Server) GetAverageSubject(data obj.ContainerSubject, reply *string) error {
	_, err := svr.Subjects[data.Subject]
	if !err {
		return errors.New("the subject has not been registered")
	}

	var average float64 = 0
	for _, grade := range svr.Subjects[data.Subject] {
		average += grade
	}
	average /= float64(len(svr.Subjects[data.Subject]))

	*reply = fmt.Sprint(average)
	return nil
}

func (svr *Server) GetAverageGeneral(data obj.ContainerEmpty, reply *string) error {
	err := len(svr.Students)
	if err == 0 {
		return errors.New("a student has never been registered")
	}

	var grades []float64
	var average float64

	for _, subjects := range svr.Students {
		average = 0
		for _, grade := range subjects {
			average += grade
		}
		average /= float64(len(subjects))
		grades = append(grades, average)
	}

	average = 0
	for _, gpa := range grades {
		average += gpa
	}
	average /= float64(len(grades))

	*reply = fmt.Sprint(average)
	return nil
}
