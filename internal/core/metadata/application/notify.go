//
// Copyright (C) 2021-2023 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"net/http"

	"github.com/edgexfoundry/go-mod-bootstrap/v3/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v3/di"
	clients "github.com/edgexfoundry/go-mod-core-contracts/v3/clients/http"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/dtos/requests"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/models"
)

func newDeviceServiceCallbackClient(ctx context.Context, dic *di.Container, deviceServiceName string) (interfaces.DeviceServiceCallbackClient, errors.EdgeX) {
	ds, err := DeviceServiceByName(deviceServiceName, ctx, dic)
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}
	return clients.NewDeviceServiceCallbackClient(ds.BaseAddress), nil
}

// validateDeviceCallback invoke device service's validation function for validating new or updated device
func validateDeviceCallback(ctx context.Context, dic *di.Container, device dtos.Device) errors.EdgeX {
	lc := container.LoggingClientFrom(dic.Get)
	deviceServiceCallbackClient, err := newDeviceServiceCallbackClient(ctx, dic, device.ServiceName)
	if err != nil {
		lc.Errorf("fail to create a device service callback client by serviceName %s, err: %v", device.ServiceName, err)
		return err
	}

	// reusing AddDeviceRequest here as it contains the protocols field and opens up
	// to other validation beyond protocols if ever needed
	request := requests.NewAddDeviceRequest(device)
	_, err = deviceServiceCallbackClient.ValidateDeviceCallback(ctx, request)
	if err != nil {
		// TODO: reconsider the validity in v3
		// allow this case for the backward-compatability in v2
		if err.Code() == http.StatusServiceUnavailable {
			lc.Warnf("Skipping device validation for device %s (device service %s unavailable)", device.Name, device.ServiceName)
		} else if err.Code() == http.StatusNotFound {
			lc.Warnf("Skipping device validation for device %s (device service %s < v2.2)", device.Name, device.ServiceName)
		} else {
			return errors.NewCommonEdgeX(errors.KindServerError, "device validation failed", err)
		}
	}

	return nil
}

// updateDeviceProfileCallback invoke device service's callback function for updating device profile
func updateDeviceProfileCallback(ctx context.Context, dic *di.Container, deviceProfile dtos.DeviceProfile) {
	lc := container.LoggingClientFrom(dic.Get)
	devices, _, err := DevicesByProfileName(0, -1, deviceProfile.Name, dic)
	if err != nil {
		lc.Errorf("fail to query associated devices by deviceProfile name %s, err: %v", deviceProfile.Name, err)
		return
	}
	// Invoke callback for each device service
	dsMap := make(map[string]bool)
	for _, d := range devices {
		if _, ok := dsMap[d.ServiceName]; ok {
			// skip the invoked device service
			continue
		}
		dsMap[d.ServiceName] = true

		deviceServiceCallbackClient, err := newDeviceServiceCallbackClient(ctx, dic, d.ServiceName)
		if err != nil {
			lc.Errorf("fail to new a device service callback client by serviceName %s, err: %v", d.ServiceName, err)
			continue
		}

		request := requests.NewDeviceProfileRequest(deviceProfile)
		response, err := deviceServiceCallbackClient.UpdateDeviceProfileCallback(ctx, request)
		if err != nil {
			lc.Errorf("fail to invoke device service callback for updating device profile %s, err: %v", deviceProfile.Name, err)
			continue
		}
		if response.StatusCode != http.StatusOK {
			lc.Errorf("fail to invoke device service callback for updating device profile %s, err: %s", deviceProfile.Name, response.Message)
		}
	}
}

// addProvisionWatcherCallback invoke device service's callback function for adding new provision watcher
func addProvisionWatcherCallback(ctx context.Context, dic *di.Container, pw dtos.ProvisionWatcher) {
	lc := container.LoggingClientFrom(dic.Get)
	deviceServiceCallbackClient, err := newDeviceServiceCallbackClient(ctx, dic, pw.ServiceName)
	if err != nil {
		lc.Errorf("fail to new a device service callback client by serviceName %s, err: %v", pw.ServiceName, err)
		return
	}

	request := requests.NewAddProvisionWatcherRequest(pw)
	response, err := deviceServiceCallbackClient.AddProvisionWatcherCallback(ctx, request)
	if err != nil {
		lc.Errorf("fail to invoke device service callback for adding  provision watcher %s, err: %v", pw.Name, err)
		return
	}
	if response.StatusCode != http.StatusOK {
		lc.Errorf("fail to invoke device service callback for adding  provision watcher %s, err: %s", pw.Name, response.Message)
	}
}

// updateProvisionWatcherCallback invoke device service's callback function for updating provision watcher
func updateProvisionWatcherCallback(ctx context.Context, dic *di.Container, serviceName string, pw models.ProvisionWatcher) {
	lc := container.LoggingClientFrom(dic.Get)
	deviceServiceCallbackClient, err := newDeviceServiceCallbackClient(ctx, dic, serviceName)
	if err != nil {
		lc.Errorf("fail to new a device service callback client by serviceName %s, err: %v", serviceName, err)
		return
	}

	request := requests.NewUpdateProvisionWatcherRequest(dtos.FromProvisionWatcherModelToUpdateDTO(pw))
	response, err := deviceServiceCallbackClient.UpdateProvisionWatcherCallback(ctx, request)
	if err != nil {
		lc.Errorf("fail to invoke device service callback for updating provision watcher %s, err: %v", pw.Name, err)
		return
	}
	if response.StatusCode != http.StatusOK {
		lc.Errorf("fail to invoke device service callback for updating provision watcher %s, err: %s", pw.Name, response.Message)
	}
}

// deleteProvisionWatcherCallback invoke device service's callback function for deleting provision watcher
func deleteProvisionWatcherCallback(ctx context.Context, dic *di.Container, pw models.ProvisionWatcher) {
	lc := container.LoggingClientFrom(dic.Get)
	deviceServiceCallbackClient, err := newDeviceServiceCallbackClient(ctx, dic, pw.ServiceName)
	if err != nil {
		lc.Errorf("fail to new a device service callback client by serviceName %s, err: %v", pw.ServiceName, err)
		return
	}
	response, err := deviceServiceCallbackClient.DeleteProvisionWatcherCallback(ctx, pw.Name)
	if err != nil {
		lc.Errorf("fail to invoke device service callback for deleting provision watcher %s, err: %v", pw.Name, err)
		return
	}
	if response.StatusCode != http.StatusOK {
		lc.Errorf("fail to invoke device service callback for deleting provision watcher %s, err: %s", pw.Name, response.Message)
	}
}

// updateDeviceServiceCallback invoke device service's callback function for updating device service
func updateDeviceServiceCallback(ctx context.Context, dic *di.Container, ds models.DeviceService) {
	lc := container.LoggingClientFrom(dic.Get)
	deviceServiceCallbackClient, err := newDeviceServiceCallbackClient(ctx, dic, ds.Name)
	if err != nil {
		lc.Errorf("fail to new a device service callback client by serviceName %s, err: %v", ds.Name, err)
		return
	}

	request := requests.NewUpdateDeviceServiceRequest(dtos.FromDeviceServiceModelToUpdateDTO(ds))
	response, err := deviceServiceCallbackClient.UpdateDeviceServiceCallback(ctx, request)
	if err != nil {
		lc.Errorf("fail to invoke device service callback for updating device service %s, err: %v", ds.Name, err)
		return
	}
	if response.StatusCode != http.StatusOK {
		lc.Errorf("fail to invoke device service callback for updating device service %s, err: %s", ds.Name, response.Message)
	}
}
