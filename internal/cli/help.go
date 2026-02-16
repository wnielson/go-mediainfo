package cli

import (
	"fmt"
	"io"
)

func Help(program string, stdout io.Writer) {
	Version(stdout)
	fmt.Fprintf(stdout, "Usage: \"%s [-Options...] FileName1 [Filename2...]\"\n", program)
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Options:")
	fmt.Fprintln(stdout, "--Help, -h")
	fmt.Fprintln(stdout, "                    Display this help and exit")
	fmt.Fprintln(stdout, "--Version")
	fmt.Fprintln(stdout, "                    Display version information and exit")
	fmt.Fprintln(stdout, "--Help-Output")
	fmt.Fprintln(stdout, "                    Display help for Output= option (templates are not implemented)")
	fmt.Fprintln(stdout, "--Help-AnOption")
	fmt.Fprintln(stdout, "                    Display help for \"AnOption\" (not implemented)")
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "--Full, -f")
	fmt.Fprintln(stdout, "                    Reserved (currently no-op)")
	fmt.Fprintln(stdout, "--Output=TEXT|JSON|XML|OLDXML|HTML|CSV|EBUCore|EBUCore_JSON|PBCore|PBCore2|Graph_Svg|Graph_Dot")
	fmt.Fprintln(stdout, "                    Select output format (templates are not implemented)")
	fmt.Fprintln(stdout, "--Language=raw")
	fmt.Fprintln(stdout, "                    Display non-translated unique identifiers (recommended for parity)")
	fmt.Fprintln(stdout, "-lang=raw")
	fmt.Fprintln(stdout, "                    Alias for --Language=raw")
	fmt.Fprintln(stdout, "--LogFile=...")
	fmt.Fprintln(stdout, "                    Save the output in the specified file")
	fmt.Fprintln(stdout, "--BOM")
	fmt.Fprintln(stdout, "                    Byte order mark for UTF-8 output (Windows only)")
	fmt.Fprintln(stdout, "--ParseSpeed=0..1")
	fmt.Fprintln(stdout, "                    Analysis speed/accuracy tradeoff (default 0.5; parity: 0.5)")
	fmt.Fprintln(stdout, "--File_TestContinuousFileNames=0|1")
	fmt.Fprintln(stdout, "                    Enable MediaInfo-style \"continuous filenames\" probing (default 0)")
	fmt.Fprintln(stdout, "--Info-Parameters")
	fmt.Fprintln(stdout, "                    Display list of inform= parameters")
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Commands:")
	fmt.Fprintln(stdout, "completion           Generate the autocompletion script for the specified shell")
	fmt.Fprintln(stdout, "help                 Help about any command")
	fmt.Fprintln(stdout, "version              Print go-mediainfo version information")
	fmt.Fprintln(stdout, "update               Update mediainfo to latest version (release builds only)")
}

func HelpNothing(program string, stdout io.Writer) {
	fmt.Fprintf(stdout, "Usage: \"%s [-Options...] FileName1 [Filename2...]\"\n", program)
	fmt.Fprintf(stdout, "\"%s --help\" for displaying more information\n", program)
}

func HelpOutput(program string, stdout io.Writer) {
	fmt.Fprintln(stdout, "--Output=...  Select an output format")
	fmt.Fprintf(stdout, "Usage: \"%s --Output=JSON FileName\"\n", program)
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Supported formats:")
	fmt.Fprintln(stdout, "TEXT, JSON, XML, OLDXML, HTML, CSV, EBUCore, EBUCore_JSON, PBCore, PBCore2, Graph_Svg, Graph_Dot")
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Note: MediaInfo output templates (e.g. \"Video;%AspectRatio%\" or \"file://...\") are not implemented.")
}

func Usage(program string, stdout io.Writer) int {
	HelpNothing(program, stdout)
	return exitError
}
