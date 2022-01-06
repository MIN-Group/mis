package Message

//go:generate msgp
type ZeroMsg struct {
	Information string `msg:"info"`
}

//将包含消息头的消息转换为byte数组
func (msg ZeroMsg) ToByteArray() (b []byte, err error) {
	b, err = msg.MarshalMsg(nil)

	/*
		header2:=MessageHeader{}
		body2:=ZeroMsg{}
		fmt.Println("b=",b)

		c,err:=body2.UnmarshalMsg(b)
		fmt.Println("c=",c)
		c,err=header2.UnmarshalMsg(c)
		fmt.Println("c=",c)

		fmt.Println("header2=",header2)
		fmt.Println("body2=",body2)
	*/

	return
}
