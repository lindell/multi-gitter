package namedflag

import (
	"fmt"
	"io"
	"strings"

	"github.com/moby/term"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	usageFmt = "Usage:\n  %s\n"
)

// NamedFlagSets stores named flag sets in the order of calling FlagSet.
type NamedFlagSets struct {
	Cmd *cobra.Command
	// Order is an ordered list of flag set names.
	Order []string
	// FlagSets stores the flag sets by name.
	FlagSets map[string]Set
}

type Set struct {
	*pflag.FlagSet
	completions map[string]CompletionFunc
}

func New(cmd *cobra.Command) NamedFlagSets {
	return NamedFlagSets{
		Cmd: cmd,
	}
}

// FlagSet returns the flag set with the given name and adds it to the
// ordered name list if it is not in there yet.
func (nfs *NamedFlagSets) FlagSet(name string) Set {
	if nfs.FlagSets == nil {
		nfs.FlagSets = map[string]Set{}
	}
	if _, ok := nfs.FlagSets[name]; !ok {
		flagSet := pflag.NewFlagSet(name, pflag.ExitOnError)
		flagSet.SortFlags = false
		nfs.FlagSets[name] = Set{
			FlagSet:     flagSet,
			completions: map[string]CompletionFunc{},
		}
		nfs.Order = append(nfs.Order, name)
	}
	return nfs.FlagSets[name]
}

// PrintSections prints the given names flag sets in sections, with the maximal given column number.
// If cols is zero, lines are not wrapped.
func PrintSections(w io.Writer, cmd *cobra.Command, fss NamedFlagSets, cols int) {
	for _, name := range fss.Order {
		fs := fss.FlagSets[name]
		if !fs.HasFlags() {
			continue
		}

		fmt.Fprintf(w, "\n%s flags:\n\n%s", strings.ToUpper(name[:1])+name[1:], fs.FlagUsagesWrapped(cols))
	}

	flagsInNamedFlags := map[string]struct{}{}
	for _, fs := range fss.FlagSets {
		fs.VisitAll(func(f *pflag.Flag) {
			flagsInNamedFlags[f.Name] = struct{}{}
		})
	}

	otherFlags := pflag.NewFlagSet("other", pflag.ExitOnError)
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if _, ok := flagsInNamedFlags[f.Name]; !ok {
			otherFlags.AddFlag(f)
		}
	})
	if otherFlags.HasFlags() {
		fmt.Fprintf(w, "\nOther flags:\n\n%s", otherFlags.FlagUsagesWrapped(cols))
	}
}

// SetUsageAndHelpFunc set both usage and help function.
// Print the flag sets we need instead of all of them.
func SetUsageAndHelpFunc(cmd *cobra.Command, fss NamedFlagSets) {
	// Add named flagsets
	for _, fs := range fss.FlagSets {
		cmd.Flags().AddFlagSet(fs.FlagSet)
		for name, comp := range fs.completions {
			_ = cmd.RegisterFlagCompletionFunc(name, comp)
		}
	}

	cols, _, _ := terminalSize(cmd.OutOrStdout())

	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		PrintSections(cmd.OutOrStderr(), cmd, fss, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		PrintSections(cmd.OutOrStdout(), cmd, fss, cols)
	})
}

// terminalSize returns the current width and height of the user's terminal. If it isn't a terminal,
// nil is returned. On error, zero values are returned for width and height.
// Usually w must be the stdout of the process. Stderr won't work.
func terminalSize(w io.Writer) (int, int, error) {
	outFd, isTerminal := term.GetFdInfo(w)
	if !isTerminal {
		return 0, 0, fmt.Errorf("given writer is no terminal")
	}
	winsize, err := term.GetWinsize(outFd)
	if err != nil {
		return 0, 0, err
	}
	return int(winsize.Width), int(winsize.Height), nil
}

type CompletionFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

func (s *Set) RegisterFlagCompletionFunc(flagName string, f CompletionFunc) error {
	s.completions[flagName] = f
	return nil
}
