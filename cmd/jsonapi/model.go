// Package jsonapi provides the API for JSON export of AHA data.
package jsonapi

import "time"

// codebeat:disable[TOO_MANY_IVARS]

// DeviceList wraps a collection of devices.
type DeviceList struct {
	NumberOfItems int      `json:"numberOfItems"` // Number of items.
	Devices       []Device `json:"devices"`       // The temperature to switch to. Same unit convention as in Thermostat.Measured.
}

// Device can represent any AHA managed hardware.
type Device struct {
	ID           string        `json:"id,omitempty"`           // A unique ID like AIN, MAC address, etc.
	InternalID   string        `json:"internalId,omitempty"`   // Internal device ID of the FRITZ!Box.
	Name         string        `json:"name,omitempty"`         // The name of the device. Can be assigned in the web gui of the FRITZ!Box.
	Properties   *Properties   `json:"properties,omitempty"`   // Static or pseudo-static properties.
	Measurements *Measurements `json:"measurements,omitempty"` // Sensor data.
	State        *State        `json:"state,omitempty"`        // State of the home.
}

// Properties refers to static or rarely changing information on the device.
type Properties struct {
	Vendor   *Vendor  `json:"vendor,omitempty"`   // Vendor data.
	Lock     *Lock    `json:"lock,omitempty"`     // Lock state.
	Warnings []string `json:"warnings,omitempty"` // Error messages, warnings etc.
}

// Vendor contains device metadata, see embedded fields.
type Vendor struct {
	Manufacturer    string `json:"manufacturer"`    // Manufacturer of the device.
	ProductName     string `json:"productName"`     // Name of the product, empty for unknown or undefined devices.
	FirmwareVersion string `json:"firmwareVersion"` // Firmware version of the device.
}

// Lock contains the lock information of a device. A "lock" prevents manual changes to the device.
type Lock struct {
	HwLock string `json:"hwLock,omitempty"` // Lock set directly on the device.
	SwLock string `json:"swLock,omitempty"` // Device locked by FRITZ!Box.
}

// Measurements indicate runtime data obtained by device senors.
type Measurements struct {
	Temperature       string     `json:"temperature,omitempty"`       // Temperature measured in °C.
	PowerConsumption  string     `json:"powerConsumption,omitempty"`  // Current power in W.
	EnergyConsumption string     `json:"energyConsumption,omitempty"` // Absolute energy consumption in Wh since the device started operating.
	AlertSignal       string     `json:"alertSignal,omitempty"`       // "ON", "OFF" (if the device reports an alert) or "" (if unknown or does not apply).
	ButtonLastPressed *time.Time `json:"buttonLastPressed,omitempty"` // Timestamp when the button was last pressed.
}

// State contains the core domain of the device, e.g. "is the switch on?", "what is the room temperature supposed to be?".
type State struct {
	Connected          bool                `json:"connected"`                    // Device connected?
	Switch             string              `json:"switch,omitempty"`             // "ON" or "OFF" or "" (if it does not apply).
	TemperatureControl *TemperatureControl `json:"temperatureControl,omitempty"` // Applies to thermostats.
	BatteryState       string              `json:"batteryState,omitempty"`       // Describes the state of the battery (if any), "OK", "LOW" or "" (if unknown).
	BatteryChargeLevel string              `json:"batteryChargeLevel,omitempty"` // Charge level of the battery (if any), ranges from 0 to 1.
}

// TemperatureControl applies to AHA devices capable of adjusting room temperature.
type TemperatureControl struct {
	Goal               string      `json:"goal,omitempty"`               // Desired temperature, user controlled.
	Saving             string      `json:"saving,omitempty"`             // Energy saving temperature.
	Comfort            string      `json:"comfort,omitempty"`            // Comfortable temperature.
	NextChange         *NextChange `json:"nextChange,omitempty"`         // Comfortable temperature.
	BatteryLow         string      `json:"batteryLow,omitempty"`         // "0" if the battery is OK, "1" if it is running low on capacity.
	BatteryChargeLevel string      `json:"batteryPercent,omitempty"`     // Battery charging percentage
	Window             string      `json:"window,omitempty"`             // "OPEN", "CLOSED" or "" (if unknown).
	WindowOpenEnd      string      `json:"windowOpenEndtime,omitempty"`  // Scheduled end of window-open state (seconds since 1970)
	Boost              bool        `json:"boost,omitempty"`              // true if boost mode is active, false if not.
	BoostEnd           string      `json:"boostActiveEndtime,omitempty"` // Scheduled end of boost time (seconds since 1970)
	Holiday            bool        `json:"holidayactive"`                // true if device is in holiday-mode, false if not.
	Summer             bool        `json:"summeractive"`                 // true if device is in summer mode (heating off), false if not.

}

// NextChange indicates the upcoming scheduled temperature change.
type NextChange struct {
	At   string `json:"at"`   // Timestamp  when the next temperature switch is scheduled. Formatted as RFC3339.
	Goal string `json:"goal"` // The temperature to switch to.
}

// codebeat:enable[TOO_MANY_IVARS]
