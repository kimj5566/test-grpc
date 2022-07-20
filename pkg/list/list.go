// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package client is a small tool to query the streaming gRPC endpoint.
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming/pkg/api/v1"
)

func List() {
	var opts []grpc.DialOption

	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))

	conn, _ := grpc.Dial("grpc-server-streaming-yv5evqvhma-uc.a.run.app:443 ", opts...)

	defer conn.Close()
	client := pb.NewTimeServiceClient(conn)

	if err := streamTime(client, 5); err != nil {
		log.Fatal(err)
	}
}

func streamTime(client pb.TimeServiceClient, duration uint) error {
	ctx := context.Background()

	resp, err := client.StreamTime(ctx, &pb.Request{
		DurationSecs: uint32(duration)})
	if err != nil {
		return fmt.Errorf("StreamTime rpc failed: %w", err)
	}
	log.Print("rpc established to timeserver, starting to stream")

	for {
		msg, err := resp.Recv()
		if err == io.EOF {
			log.Printf("end of stream")
			return nil
		} else if err != nil {
			return fmt.Errorf("error receiving message: %w", err)
		}

		ts, err := ptypes.Timestamp(msg.GetCurrentTime())
		if err != nil {
			return fmt.Errorf("failed to parse timestamp %v: %w", msg.GetCurrentTime(), err)
		}
		log.Printf("received message: current_timestamp: %v", ts.Format(time.RFC3339))
	}
}
