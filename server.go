package main

import (
	"net/http"

	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/docker/go-plugins-helpers/volume"
)

func newVolumeServer(driver volume.Driver) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/Plugin.Activate", func(w http.ResponseWriter, _ *http.Request) {
		sdk.EncodeResponse(w, map[string][]string{"Implements": {"VolumeDriver"}}, false)
	})
	mux.HandleFunc("/VolumeDriver.Create", func(w http.ResponseWriter, r *http.Request) {
		var req volume.CreateRequest
		if sdk.DecodeRequest(w, r, &req) == nil {
			writeVolumeResponse(w, struct{}{}, driver.Create(&req))
		}
	})
	mux.HandleFunc("/VolumeDriver.Get", func(w http.ResponseWriter, r *http.Request) {
		var req volume.GetRequest
		if sdk.DecodeRequest(w, r, &req) == nil {
			resp, err := driver.Get(&req)
			writeVolumeResponse(w, resp, err)
		}
	})
	mux.HandleFunc("/VolumeDriver.List", func(w http.ResponseWriter, _ *http.Request) {
		resp, err := driver.List()
		writeVolumeResponse(w, resp, err)
	})
	mux.HandleFunc("/VolumeDriver.Remove", func(w http.ResponseWriter, r *http.Request) {
		var req volume.RemoveRequest
		if sdk.DecodeRequest(w, r, &req) == nil {
			writeVolumeResponse(w, struct{}{}, driver.Remove(&req))
		}
	})
	mux.HandleFunc("/VolumeDriver.Path", func(w http.ResponseWriter, r *http.Request) {
		var req volume.PathRequest
		if sdk.DecodeRequest(w, r, &req) == nil {
			resp, err := driver.Path(&req)
			writeVolumeResponse(w, resp, err)
		}
	})
	mux.HandleFunc("/VolumeDriver.Mount", func(w http.ResponseWriter, r *http.Request) {
		var req volume.MountRequest
		if sdk.DecodeRequest(w, r, &req) == nil {
			resp, err := driver.Mount(&req)
			writeVolumeResponse(w, resp, err)
		}
	})
	mux.HandleFunc("/VolumeDriver.Unmount", func(w http.ResponseWriter, r *http.Request) {
		var req volume.UnmountRequest
		if sdk.DecodeRequest(w, r, &req) == nil {
			writeVolumeResponse(w, struct{}{}, driver.Unmount(&req))
		}
	})
	mux.HandleFunc("/VolumeDriver.Capabilities", func(w http.ResponseWriter, _ *http.Request) {
		sdk.EncodeResponse(w, driver.Capabilities(), false)
	})

	server := &http.Server{Handler: mux}
	// Finish each local RPC without retaining an idle connection to Docker.
	server.SetKeepAlivesEnabled(false)
	return server
}

func writeVolumeResponse(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		sdk.EncodeResponse(w, volume.NewErrorResponse(err.Error()), true)
		return
	}
	sdk.EncodeResponse(w, resp, false)
}
