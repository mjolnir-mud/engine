package engine

//func TestAddCommand(t *testing.T) {
//	Init("test", []Plugin{&testPlugin{}})
//
//	AddCommand(&cobra.Command{
//		Use: "test",
//		Run: func(cmd *cobra.Command, args []string) {
//			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "testing")
//		},
//	})
//
//	b := bytes.NewBufferString("")
//	state.baseCommand.SetOut(b)
//
//	state.baseCommand.SetArgs([]string{"test"})
//	ExecCommand()
//
//	out, err := ioutil.ReadAll(b)
//
//	if err != nil {
//		t.Fatalf("error reading output: %s", err)
//	}
//
//	assert.Equal(t, "testing", string(out))
//}
