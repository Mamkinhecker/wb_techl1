package main

type OldUSB interface {
	ConnectWithUSB() string
}

type NewUSBC interface {
	ConnectWithUSBC() string
}

type USBDevice struct {
	name string
}

func (d *USBDevice) ConnectWithUSB() string {
	return d.name + " подключено через USB‑A"
}

type Computer struct{}

func (c *Computer) Connect(device NewUSBC) {
	result := device.ConnectWithUSBC()
	println(result)
}

type Adapter struct {
	oldDevice *USBDevice
}

func (a *Adapter) ConnectWithUSBC() string {
	return a.oldDevice.ConnectWithUSB() + " (через адаптер USB‑A → USB‑C)"
}

func main() {
	oldDevice := &USBDevice{name: "Флешка"}

	adapter := &Adapter{oldDevice: oldDevice}

	computer := &Computer{}
	computer.Connect(adapter)
}
