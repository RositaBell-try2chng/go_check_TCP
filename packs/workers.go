package portscan

import (
	"context"
)

func generator(
	ctx context.Context,
	hosts []string,
	ports []uint16,
	jobs chan<- Job,
	results chan<- Result,
	) {
		exists := make(map[string]struct{})
		for _, host := range hosts {
			if len(host) == 0 { continue }
			ips, err := resolveHost(host)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case results <- Result{
					Job: Job{
						Host: host,
						IP: nil,
						Port: 0,
					},
					State: StateError,
					Duration: 0,
					Err: err,
					}:
				}
				continue
			}

			for _, ip := range ips {
				ipStr := ip.String()
				if _, ok := exists[ipStr]; ok {
					continue
				}
				exists[ipStr] = struct{}{}
				for _, port := range ports {
					select {
					case <-ctx.Done():
						return
					case jobs <- Job{
						Host: host,
						IP: ip,
						Port: port,
						}:
					}
				}
			}
		}
}

func (s *Scanner)worker(
	ctx context.Context,
	jobs <-chan Job,
	results chan<- Result,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			select {
			case <-ctx.Done():
				return
			case results <- s.doIt(ctx, job):
			}
		}
	}
}

func collector(
	ctx context.Context,
	results <-chan Result,
) []Result {
	collected := make([]Result, 0)
	for {
		select {
		case <- ctx.Done():
			return collected
		case res, ok := <-results:
			if !ok {
				return collected
			}
			collected = append(collected, res)
		}
	}
}