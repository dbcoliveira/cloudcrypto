package main

import (
	"fmt"
	"moonraker/api"
	"os"
)

func main() {
	var cloudcrypt api.Login

	cloudcrypt.Connect(("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODEwMDY0MDIsInVzZXIiOiJ0ZXN0QHRlc3QuY29tIn0.cm6eDFvahjJBUBKSZfNBJK5Zf8hSLKzHBinVzGU9F_c"))

	sink := api.NewSink(os.Args[1])
	cloudcrypt.AssignSink(sink)

	resEnc, err := cloudcrypt.Encrypt(os.Args[2])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resEnc.Data.Text)

	resDec, err := cloudcrypt.Decrypt(resEnc.Data.Text)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resDec.Data.Text)
}
