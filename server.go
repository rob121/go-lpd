package lpd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

// QueueHandler is an interface like http handler
type QueueHandler interface {
	ProcessJob(PrintJob)
}

// QueueFunc add ProcessJob method to a function
type QueueFunc func(PrintJob)

// ProcessJob calls qh(p)
func (qh QueueFunc) ProcessJob(p PrintJob) {
	qh(p)
}

// A Server listen and store Jobs
type Server struct {
	port   int
	Jobs   []PrintJob
	queues map[string]QueueHandler
	// lock sync.Mutex
}

// Handle adds an handler to a queue
func (s *Server) Handle(QueueName string, q QueueHandler) {
	if s.queues == nil {
		s.queues = make(map[string]QueueHandler)
	}
	s.queues[QueueName] = q
}

// HandleFunc registers the handler function for the given pattern.
func (s *Server) HandleFunc(q string, f func(p PrintJob)) {
	s.Handle(q, QueueFunc(f))
}

// NewServer creates and returns an instance of Server that
// will listen at address. The standard port
// for an LPD server is 515.
func NewServer(port int) (*Server, error) {
	srv := new(Server)

	srv.port = port
	return srv, nil
}

// Serve our requests
func (s *Server) Serve() {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(s.port))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	// fmt.Println("Listening on 0.0.0.0:5515")

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go s.processClient(conn)
	}
}

func (s *Server) processClient(conn net.Conn) {
	reader := bufio.NewReader(conn)
	defer conn.Close()
	// writer := bufio.NewWriter(conn)
	data, _ := reader.ReadBytes(0x0a)
	// fmt.Printf("Data recieved : %v\n", data)

	cmd, err := unmarshalCommand(data[:len(data)-1])
	
	log.Println(cmd)
	
	if err != nil {
		log.Fatal("Can't unmarshal")
	}

	if cmd.Code != ReceivePrintJob {
		fmt.Println("Not Implemented")
		return
	}

	if queue := s.queues[cmd.Queue]; queue == nil {
		nackSubCommand(conn)
		fmt.Printf("unknown queue '%s' in '%v'\n", cmd.Queue, s.queues)
		return
	}

	ackSubCommand(conn)
	// fmt.Printf("Bytes : %v Ok next...", b)

	job, err := receiveJob(conn, conn)
	if err != nil {
		log.Println(err)
		return
	}

	job.QueueName = cmd.Queue
	
	// fmt.Printf("Job Recieved : %v", job)
	s.queues[cmd.Queue].ProcessJob(*job)

	return
}
