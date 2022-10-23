package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"modular-netstalking/lib"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	flag.StringVar(&lib.DictPath, "d", "./assets/data/rtsp-paths.txt", "dictionary path")
}

const _RTSP_TPL = "%s %s RTSP/1.0\r\n" +
	"Accept: application/sdp\r\n" +
	"CSeq: %d\r\n" +
	"User-Agent: Lavf59.16.100\r\n\r\n"

type RTSPConn struct {
	Hi   *lib.HostInfo
	conn net.Conn
	cseq int
}

func (r *RTSPConn) Request(req string) (code int, err error) {
	_ = r.conn.SetDeadline(time.Now().Add(2 * time.Second))

	if _, err = r.conn.Write([]byte(req)); err != nil {
		return
	}

	data := make([]byte, 1024)
	n, err := r.conn.Read(data)
	if err != nil {
		return
	}

	f := strings.Fields(string(data[:n]))
	if len(f) >= 2 && strings.HasPrefix(f[0], "RTSP/") {
		return strconv.Atoi(f[1])
	}

	return 0, errors.New("Bad response")
}

func (r *RTSPConn) Query(path string) string {
	r.cseq++

	return fmt.Sprintf(_RTSP_TPL, "DESCRIBE", fmt.Sprintf("rtsp://%s%s", r.Hi.Host, path), r.cseq)
}

func (rc *RTSPConn) CheckPath(path string) bool {
	code, err := rc.Request(rc.Query(path))
	if err != nil || code == 403 {
		return true
	}

	if code == 200 {
		rc.Hi.Attrs["paths"] = append((rc.Hi.Attrs["paths"]).([]string), path)
	}

	return false
}

func main() {
	dict, err := lib.FileToArray(lib.DictPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dict open: ", err)
		os.Exit(255)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ln := scanner.Text()
		hi := lib.HostInfoFromJson(ln)
		hi.Attrs["paths"] = []string{}
		c, err := net.DialTimeout("tcp", hi.Host, time.Second)
		if err != nil {
			continue
		}
		conn := RTSPConn{Hi: &hi, conn: c}
		for _, path := range dict {
			if conn.CheckPath(path) {
				break
			}
		}
		fmt.Println(hi)
	}
}
