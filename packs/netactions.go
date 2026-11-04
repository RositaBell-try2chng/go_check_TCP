package portscan

import (
	"net"
	"time"
	"errors"
	"syscall"
	"os"
	"strconv"
	"context"
)

func (s *Scanner)doIt(ctx context.Context, job Job) Result {
	address := net.JoinHostPort(job.IP.String(), strconv.Itoa(int(job.Port)))

	dialer := net.Dialer{
		Timeout: s.tm,
	}
	start := time.Now()
	conn, err := dialer.DialContext(ctx, s.pr, address)
	latency := time.Since(start)
	
	result := Result{
		Job: Job{
			Host: job.Host,
			IP: job.IP,
			Port: job.Port,
		},
		Duration: latency,
	}

	if err != nil {
		result.State = wichErr(ctx, err)
		result.Err = err
		return result
	}
	defer conn.Close()
	result.State = StateOpen
	return result

}

func resolveHost(host string) ([]net.IP, error) {
	if ip := net.ParseIP(host); ip != nil {
		return []net.IP{ip}, nil
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

func wichErr(ctx context.Context, err error) State {
	if ctx.Err() != nil {
		return StateCanceled
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		
		var dnsErr *net.DNSError
		if errors.As(opErr.Err, &dnsErr) {
			return StateDNSFailed
		}

		var syscallErr *os.SyscallError
		if errors.As(opErr.Err, &syscallErr) {
			switch syscallErr.Err {
			case syscall.ECONNREFUSED:
				return StateClosed

			case syscall.ETIMEDOUT:
				return StateTimeout

			case syscall.ENETUNREACH:
				return StateUnreachable

			case syscall.EHOSTUNREACH:
				return StateUnreachable

			case syscall.EACCES, syscall.EPERM:
				return StatePermissionDenied
			}
		}
	}
	return StateError
}