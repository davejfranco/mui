package user

import "os/exec"

/*
func main() {
	cmd := "convert"
	args := []string{"-resize", "50%", "foo.jpg", "foo.half.jpg"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Successfully halved image in size")
}
*/
func Execute(comand string, args ...string) error {

	if err := exec.Command(comand, args...).Run(); err != nil {
		return err
	}
	return nil
}
