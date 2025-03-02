// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

package feeder

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	cfg "github.com/kubearmor/KubeArmor/KubeArmor/config"
	tp "github.com/kubearmor/KubeArmor/KubeArmor/types"
)

// ======================= //
// == Security Policies == //
// ======================= //

// getProtocolFromName Function
func getProtocolFromName(proto string) string {
	switch strings.ToLower(proto) {
	case "tcp":
		return "protocol=TCP,type=SOCK_STREAM"
	case "udp":
		return "protocol=UDP,type=SOCK_DGRAM"
	case "icmp":
		return "protocol=ICMP,type=SOCK_RAW"
	case "raw":
		return "type=SOCK_RAW"
	default:
		return "unknown"
	}
}

func getFileProcessUID(path string) string {
	info, err := os.Stat(path)
	if err == nil {
		stat := info.Sys().(*syscall.Stat_t)
		uid := stat.Uid

		return strconv.Itoa(int(uid))
	}

	return ""
}

// getOperationAndCapabilityFromName Function
func getOperationAndCapabilityFromName(capName string) (op, cap string) {
	switch strings.ToLower(capName) {
	case "net_raw":
		op = "Network"
		cap = "SOCK_RAW"
	default:
		return "", "unknown"
	}

	return op, cap
}

// newMatchPolicy Function
func (fd *Feeder) newMatchPolicy(policyEnabled int, policyName, src string, mp interface{}) tp.MatchPolicy {
	match := tp.MatchPolicy{
		PolicyName: policyName,
		Source:     src,
	}

	if ppt, ok := mp.(tp.ProcessPathType); ok {
		match.Severity = strconv.Itoa(ppt.Severity)
		match.Tags = ppt.Tags
		match.Message = ppt.Message

		match.Operation = "Process"
		match.Resource = ppt.Path
		match.ResourceType = "Path"

		match.OwnerOnly = ppt.OwnerOnly

		if policyEnabled == tp.KubeArmorPolicyAudited && ppt.Action == "Allow" {
			match.Action = "Audit (" + ppt.Action + ")"
		} else if policyEnabled == tp.KubeArmorPolicyAudited && ppt.Action == "Block" {
			match.Action = "Audit (" + ppt.Action + ")"
		} else {
			match.Action = ppt.Action
		}
	} else if pdt, ok := mp.(tp.ProcessDirectoryType); ok {
		match.Severity = strconv.Itoa(pdt.Severity)
		match.Tags = pdt.Tags
		match.Message = pdt.Message

		match.Operation = "Process"
		match.Resource = pdt.Directory
		match.ResourceType = "Directory"

		match.OwnerOnly = pdt.OwnerOnly
		match.Recursive = pdt.Recursive

		if policyEnabled == tp.KubeArmorPolicyAudited && pdt.Action == "Allow" {
			match.Action = "Audit (" + pdt.Action + ")"
		} else if policyEnabled == tp.KubeArmorPolicyAudited && pdt.Action == "Block" {
			match.Action = "Audit (" + pdt.Action + ")"
		} else {
			match.Action = pdt.Action
		}
	} else if ppt, ok := mp.(tp.ProcessPatternType); ok {
		match.Severity = strconv.Itoa(ppt.Severity)
		match.Tags = ppt.Tags
		match.Message = ppt.Message

		match.Operation = "Process"
		match.Resource = ppt.Pattern
		match.ResourceType = "" // to be defined based on the pattern matching syntax

		match.OwnerOnly = ppt.OwnerOnly

		if policyEnabled == tp.KubeArmorPolicyAudited && ppt.Action == "Allow" {
			match.Action = "Audit (" + ppt.Action + ")"
		} else if policyEnabled == tp.KubeArmorPolicyAudited && ppt.Action == "Block" {
			match.Action = "Audit (" + ppt.Action + ")"
		} else {
			match.Action = ppt.Action
		}
	} else if fpt, ok := mp.(tp.FilePathType); ok {
		match.Severity = strconv.Itoa(fpt.Severity)
		match.Tags = fpt.Tags
		match.Message = fpt.Message

		match.Operation = "File"
		match.Resource = fpt.Path
		match.ResourceType = "Path"

		match.OwnerOnly = fpt.OwnerOnly
		match.ReadOnly = fpt.ReadOnly

		if policyEnabled == tp.KubeArmorPolicyAudited && fpt.Action == "Allow" {
			match.Action = "Audit (" + fpt.Action + ")"
		} else if policyEnabled == tp.KubeArmorPolicyAudited && fpt.Action == "Block" {
			match.Action = "Audit (" + fpt.Action + ")"
		} else {
			match.Action = fpt.Action
		}
	} else if fdt, ok := mp.(tp.FileDirectoryType); ok {
		match.Severity = strconv.Itoa(fdt.Severity)
		match.Tags = fdt.Tags
		match.Message = fdt.Message

		match.Operation = "File"
		match.Resource = fdt.Directory
		match.ResourceType = "Directory"

		match.OwnerOnly = fdt.OwnerOnly
		match.ReadOnly = fdt.ReadOnly
		match.Recursive = fdt.Recursive

		if policyEnabled == tp.KubeArmorPolicyAudited && fdt.Action == "Allow" {
			match.Action = "Audit (" + fdt.Action + ")"
		} else if policyEnabled == tp.KubeArmorPolicyAudited && fdt.Action == "Block" {
			match.Action = "Audit (" + fdt.Action + ")"
		} else {
			match.Action = fdt.Action
		}
	} else if fpt, ok := mp.(tp.FilePatternType); ok {
		match.Severity = strconv.Itoa(fpt.Severity)
		match.Tags = fpt.Tags
		match.Message = fpt.Message
		match.Operation = "File"
		match.Resource = fpt.Pattern
		match.ResourceType = "" // to be defined based on the pattern matching syntax

		match.OwnerOnly = fpt.OwnerOnly
		match.ReadOnly = fpt.ReadOnly

		if policyEnabled == tp.KubeArmorPolicyAudited && fpt.Action == "Allow" {
			match.Action = "Audit (" + fpt.Action + ")"
		} else if policyEnabled == tp.KubeArmorPolicyAudited && fpt.Action == "Block" {
			match.Action = "Audit (" + fpt.Action + ")"
		} else {
			match.Action = fpt.Action
		}
	} else if npt, ok := mp.(tp.NetworkProtocolType); ok {
		match.Severity = strconv.Itoa(npt.Severity)
		match.Tags = npt.Tags
		match.Message = npt.Message

		match.Operation = "Network"
		match.Resource = getProtocolFromName(npt.Protocol)
		match.ResourceType = "Protocol"

		if policyEnabled == tp.KubeArmorPolicyAudited && npt.Action == "Allow" {
			match.Action = "Audit (" + npt.Action + ")"
		} else if policyEnabled == tp.KubeArmorPolicyAudited && npt.Action == "Block" {
			match.Action = "Audit (" + npt.Action + ")"
		} else if policyEnabled == tp.KubeArmorPolicyEnabled && fd.IsGKE && npt.Action == "Block" {
			match.Action = "Audit (" + npt.Action + ")"
		} else {
			match.Action = npt.Action
		}
	} else if cct, ok := mp.(tp.CapabilitiesCapabilityType); ok {
		match.Severity = strconv.Itoa(cct.Severity)
		match.Tags = cct.Tags
		match.Message = cct.Message

		op, cap := getOperationAndCapabilityFromName(cct.Capability)

		match.Operation = op
		match.Resource = cap
		match.ResourceType = "Capability"

		if policyEnabled == tp.KubeArmorPolicyAudited && cct.Action == "Allow" {
			match.Action = "Audit (" + cct.Action + ")"
		} else if policyEnabled == tp.KubeArmorPolicyAudited && cct.Action == "Block" {
			match.Action = "Audit (" + cct.Action + ")"
		} else {
			match.Action = cct.Action
		}
	} else if smt, ok := mp.(tp.SyscallMatchType); ok {
		match.Severity = strconv.Itoa(smt.Severity)
		match.Tags = smt.Tags
		match.Message = smt.Message
		match.Operation = "Syscall"
		match.ResourceType = strings.ToUpper(smt.Syscalls[0])
		match.Action = "Audit"
	} else if smpt, ok := mp.(tp.SyscallMatchPathType); ok {
		match.Severity = strconv.Itoa(smpt.Severity)
		match.Tags = smpt.Tags
		match.Message = smpt.Message
		match.Action = "Audit"
		match.Operation = "Syscall"
		match.Resource = smpt.Path
		match.ResourceType = strings.ToUpper(smpt.Syscalls[0])

	} else {
		return tp.MatchPolicy{}
	}

	return match
}

// UpdateSecurityPolicies Function
func (fd *Feeder) UpdateSecurityPolicies(action string, endPoint tp.EndPoint) {
	name := endPoint.NamespaceName + "_" + endPoint.EndPointName

	if action == "DELETED" {
		delete(fd.SecurityPolicies, name)
		return
	}

	// ADDED | MODIFIED
	matches := tp.MatchPolicies{}

	for _, secPolicy := range endPoint.SecurityPolicies {
		policyName := secPolicy.Metadata["policyName"]

		if len(secPolicy.Spec.AppArmor) > 0 {
			continue
		}

		for _, path := range secPolicy.Spec.Process.MatchPaths {
			fromSource := ""

			if len(path.FromSource) == 0 {
				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, path)
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range path.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, path)
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		for _, dir := range secPolicy.Spec.Process.MatchDirectories {
			fromSource := ""

			if len(dir.FromSource) == 0 {
				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, dir)
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range dir.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, dir)
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		for _, patt := range secPolicy.Spec.Process.MatchPatterns {
			if len(patt.Pattern) == 0 {
				continue
			}

			fromSource := ""

			match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, patt)

			regexpComp, err := regexp.Compile(patt.Pattern)
			if err != nil {
				fd.Debugf("MatchPolicy Regexp compilation error: %s\n", patt.Pattern)
				continue
			}
			match.Regexp = regexpComp
			// Using 'Glob' despite compiling 'Regexp', since we don't have
			// a native pattern matching design yet and 'Glob' is more similar
			// to AppArmor's pattern matching syntax.
			match.ResourceType = "Glob"

			matches.Policies = append(matches.Policies, match)
		}

		for _, path := range secPolicy.Spec.File.MatchPaths {
			fromSource := ""

			if len(path.FromSource) == 0 {
				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, path)
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range path.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, path)
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		for _, dir := range secPolicy.Spec.File.MatchDirectories {
			fromSource := ""

			if len(dir.FromSource) == 0 {
				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, dir)
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range dir.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, dir)
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		for _, patt := range secPolicy.Spec.File.MatchPatterns {
			if len(patt.Pattern) == 0 {
				continue
			}

			fromSource := ""

			match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, patt)

			regexpComp, err := regexp.Compile(patt.Pattern)
			if err != nil {
				fd.Debugf("MatchPolicy Regexp compilation error: %s\n", patt.Pattern)
				continue
			}
			match.Regexp = regexpComp
			// Using 'Glob' despite compiling 'Regexp', since we don't have
			// a native pattern matching design yet and 'Glob' is more similar
			// to AppArmor's pattern matching syntax.
			match.ResourceType = "Glob"

			matches.Policies = append(matches.Policies, match)
		}

		for _, proto := range secPolicy.Spec.Network.MatchProtocols {
			if len(proto.Protocol) == 0 {
				continue
			}

			fromSource := ""

			if len(proto.FromSource) == 0 {
				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, proto)
				if len(match.Resource) == 0 {
					continue
				}
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range proto.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, proto)
				if len(match.Resource) == 0 {
					continue
				}
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}

		}

		for _, cap := range secPolicy.Spec.Capabilities.MatchCapabilities {
			if len(cap.Capability) == 0 {
				continue
			}

			fromSource := ""

			if len(cap.FromSource) == 0 {
				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, cap)
				if len(match.Resource) == 0 {
					continue
				}
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range cap.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, cap)
				if len(match.Resource) == 0 {
					continue
				}
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		// MatchSyscalls
		for _, syscallRule := range secPolicy.Spec.Syscalls.MatchSyscalls {
			if len(syscallRule.Syscalls) == 0 {
				continue
			}
			fromSource := ""
			syscall := tp.SyscallMatchType{
				Tags:     syscallRule.Tags,
				Message:  syscallRule.Message,
				Severity: syscallRule.Severity,
			}
			if len(syscallRule.FromSource) == 0 {
				for _, syscallName := range syscallRule.Syscalls {
					syscall.Syscalls = []string{syscallName}
					match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, syscall)
					if len(match.ResourceType) == 0 {
						continue
					}
					matches.Policies = append(matches.Policies, match)
				}
				continue
			}

			for _, src := range syscallRule.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else if len(src.Dir) > 0 {
					fromSource = src.Dir
					if !strings.HasSuffix(fromSource, "/") {
						fromSource += "/"
					}
				} else {
					continue
				}
				for _, syscallName := range syscallRule.Syscalls {
					syscall.Syscalls = []string{syscallName}
					match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, syscall)
					if len(match.ResourceType) == 0 {
						continue
					}
					match.IsFromSource = len(fromSource) > 0
					match.Recursive = len(src.Path) == 0 && src.Recursive
					matches.Policies = append(matches.Policies, match)
				}

			}
		}
		// SyscallsMatchPath
		for _, syscallRule := range secPolicy.Spec.Syscalls.MatchPaths {
			if len(syscallRule.Path) == 0 || len(syscallRule.Syscalls) == 0 {
				continue
			}
			fromSource := ""
			syscall := tp.SyscallMatchPathType{
				Tags:     syscallRule.Tags,
				Message:  syscallRule.Message,
				Severity: syscallRule.Severity,
				Path:     syscallRule.Path,
			}
			if len(syscallRule.FromSource) == 0 {
				for _, syscallName := range syscallRule.Syscalls {
					syscall.Syscalls = []string{syscallName}
					match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, syscall)
					if len(match.ResourceType) == 0 && len(match.Resource) == 0 {
						continue
					}
					match.ReadOnly = syscallRule.Recursive
					matches.Policies = append(matches.Policies, match)
				}
				continue
			}

			for _, src := range syscallRule.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else if len(src.Dir) > 0 {
					fromSource = src.Dir
					if !strings.HasSuffix(fromSource, "/") {
						fromSource += "/"
					}
				} else {
					continue
				}
				for _, syscallName := range syscallRule.Syscalls {
					syscall.Syscalls = []string{syscallName}
					match := fd.newMatchPolicy(endPoint.PolicyEnabled, policyName, fromSource, syscall)
					if len(match.ResourceType) == 0 && len(match.Resource) == 0 {
						continue
					}
					match.IsFromSource = len(fromSource) > 0
					match.Recursive = len(src.Path) == 0 && src.Recursive
					match.ReadOnly = syscallRule.Recursive
					matches.Policies = append(matches.Policies, match)
				}
			}

		}
	}

	fd.SecurityPoliciesLock.Lock()
	fd.SecurityPolicies[name] = matches
	fd.SecurityPoliciesLock.Unlock()
}

// ============================ //
// == Host Security Policies == //
// ============================ //

// UpdateHostSecurityPolicies Function
func (fd *Feeder) UpdateHostSecurityPolicies(action string, secPolicies []tp.HostSecurityPolicy) {
	if action == "DELETED" {
		delete(fd.SecurityPolicies, fd.Node.NodeName)
		return
	}

	// ADDED | MODIFIED
	matches := tp.MatchPolicies{}

	for _, secPolicy := range secPolicies {
		policyName := secPolicy.Metadata["policyName"]

		if len(secPolicy.Spec.AppArmor) > 0 {
			continue
		}

		for _, path := range secPolicy.Spec.Process.MatchPaths {
			fromSource := ""

			if len(path.FromSource) == 0 {
				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, path)
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range path.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, path)
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		for _, dir := range secPolicy.Spec.Process.MatchDirectories {
			fromSource := ""

			if len(dir.FromSource) == 0 {
				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, dir)
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range dir.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, dir)
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		for _, patt := range secPolicy.Spec.Process.MatchPatterns {
			if len(patt.Pattern) == 0 {
				continue
			}

			fromSource := ""

			match := fd.newMatchPolicy(tp.KubeArmorPolicyEnabled, policyName, fromSource, patt)

			regexpComp, err := regexp.Compile(patt.Pattern)
			if err != nil {
				fd.Debugf("MatchPolicy Regexp compilation error: %s\n", patt.Pattern)
				continue
			}
			match.Regexp = regexpComp
			// Using 'Glob' despite compiling 'Regexp', since we don't have
			// a native pattern matching design yet and 'Glob' is more similar
			// to AppArmor's pattern matching syntax.
			match.ResourceType = "Glob"

			matches.Policies = append(matches.Policies, match)
		}

		for _, path := range secPolicy.Spec.File.MatchPaths {
			fromSource := ""

			if len(path.FromSource) == 0 {
				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, path)
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range path.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, path)
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		for _, dir := range secPolicy.Spec.File.MatchDirectories {
			fromSource := ""

			if len(dir.FromSource) == 0 {
				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, dir)
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range dir.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, dir)
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		for _, patt := range secPolicy.Spec.File.MatchPatterns {
			if len(patt.Pattern) == 0 {
				continue
			}

			fromSource := ""

			match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, patt)

			regexpComp, err := regexp.Compile(patt.Pattern)
			if err != nil {
				fd.Debugf("MatchPolicy Regexp compilation error: %s\n", patt.Pattern)
				continue
			}
			match.Regexp = regexpComp
			// Using 'Glob' despite compiling 'Regexp', since we don't have
			// a native pattern matching design yet and 'Glob' is more similar
			// to AppArmor's pattern matching syntax.
			match.ResourceType = "Glob"

			matches.Policies = append(matches.Policies, match)
		}

		for _, proto := range secPolicy.Spec.Network.MatchProtocols {
			if len(proto.Protocol) == 0 {
				continue
			}

			fromSource := ""

			if len(proto.FromSource) == 0 {
				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, proto)
				if len(match.Resource) == 0 {
					continue
				}
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range proto.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, proto)
				if len(match.Resource) == 0 {
					continue
				}
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		for _, cap := range secPolicy.Spec.Capabilities.MatchCapabilities {
			if len(cap.Capability) == 0 {
				continue
			}

			fromSource := ""

			if len(cap.FromSource) == 0 {
				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, cap)
				if len(match.Resource) == 0 {
					continue
				}
				matches.Policies = append(matches.Policies, match)
				continue
			}

			for _, src := range cap.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else {
					continue
				}

				match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, cap)
				if len(match.Resource) == 0 {
					continue
				}
				match.IsFromSource = len(fromSource) > 0
				matches.Policies = append(matches.Policies, match)
			}
		}

		// MatchSyscalls
		for _, syscallRule := range secPolicy.Spec.Syscalls.MatchSyscalls {
			if len(syscallRule.Syscalls) == 0 {
				continue
			}
			fromSource := ""
			syscall := tp.SyscallMatchType{
				Tags:     syscallRule.Tags,
				Message:  syscallRule.Message,
				Severity: syscallRule.Severity,
			}
			if len(syscallRule.FromSource) == 0 {
				for _, syscallName := range syscallRule.Syscalls {
					syscall.Syscalls = []string{syscallName}
					match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, syscall)
					if len(match.ResourceType) == 0 {
						continue
					}
					matches.Policies = append(matches.Policies, match)
				}
				continue
			}

			for _, src := range syscallRule.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else if len(src.Dir) > 0 {
					fromSource = src.Dir
					if !strings.HasSuffix(fromSource, "/") {
						fromSource += "/"
					}
				} else {
					continue
				}
				for _, syscallName := range syscallRule.Syscalls {
					syscall.Syscalls = []string{syscallName}
					match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, syscall)
					if len(match.ResourceType) == 0 {
						continue
					}
					match.IsFromSource = len(fromSource) > 0
					match.Recursive = len(src.Path) == 0 && src.Recursive
					matches.Policies = append(matches.Policies, match)
				}

			}
		}
		// SyscallsMatchPath
		for _, syscallRule := range secPolicy.Spec.Syscalls.MatchPaths {
			if len(syscallRule.Path) == 0 || len(syscallRule.Syscalls) == 0 {
				continue
			}
			fromSource := ""
			syscall := tp.SyscallMatchPathType{
				Tags:     syscallRule.Tags,
				Message:  syscallRule.Message,
				Severity: syscallRule.Severity,
			}
			if len(syscallRule.FromSource) == 0 {
				for _, syscallName := range syscallRule.Syscalls {
					syscall.Syscalls = []string{syscallName}
					match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, syscall)
					if len(match.ResourceType) == 0 && len(match.Resource) == 0 {
						continue
					}
					matches.Policies = append(matches.Policies, match)
					match.Source = syscallRule.Path
				}
				continue
			}

			for _, src := range syscallRule.FromSource {
				if len(src.Path) > 0 {
					fromSource = src.Path
				} else if len(src.Dir) > 0 {
					fromSource = src.Dir
					if !strings.HasSuffix(fromSource, "/") {
						fromSource += "/"
					}
				} else {
					continue
				}
				for _, syscallName := range syscallRule.Syscalls {
					syscall.Syscalls = []string{syscallName}
					match := fd.newMatchPolicy(fd.Node.PolicyEnabled, policyName, fromSource, syscall)
					if len(match.ResourceType) == 0 && len(match.Resource) == 0 {
						continue
					}
					match.IsFromSource = len(fromSource) > 0
					match.Recursive = len(src.Path) == 0 && src.Recursive
					matches.Policies = append(matches.Policies, match)
				}
			}

		}
	}

	fd.SecurityPoliciesLock.Lock()
	fd.SecurityPolicies[fd.Node.NodeName] = matches
	fd.SecurityPoliciesLock.Unlock()
}

// ===================== //
// == Default Posture == //
// ===================== //

// UpdateDefaultPosture Function
func (fd *Feeder) UpdateDefaultPosture(action string, namespace string, defaultPosture tp.DefaultPosture) {

	fd.DefaultPosturesLock.Lock()
	defer fd.DefaultPosturesLock.Unlock()

	if action == "DELETED" {
		delete(fd.DefaultPostures, namespace)
	} else { // ADDED or MODIFIED
		fd.DefaultPostures[namespace] = defaultPosture
	}
}

// Update Log Fields based on default posture and visibility configuration and return false if no updates
func setLogFields(log *tp.Log, existAllowPolicy bool, defaultPosture string, visibility, containerEvent bool) bool {
	if existAllowPolicy && defaultPosture == "audit" && (*log).Result == "Passed" {
		if containerEvent {
			(*log).Type = "MatchedPolicy"
		} else {
			(*log).Type = "MatchedHostPolicy"
		}

		(*log).PolicyName = "DefaultPosture"
		(*log).Enforcer = "eBPF Monitor"
		(*log).Action = "Audit"

		return true
	}

	if visibility {
		if containerEvent {
			(*log).Type = "ContainerLog"
		} else {
			(*log).Type = "HostLog"
		}

		return true
	}

	return false
}

// ==================== //
// == Policy Matches == //
// ==================== //

func getDirectoryPart(path string) string {
	dir := filepath.Dir(path)
	if strings.HasPrefix(dir, "/") {
		return dir + "/"
	}
	return "__not_absolute_path__"
}

// UpdateMatchedPolicy Function
func (fd *Feeder) UpdateMatchedPolicy(log tp.Log) tp.Log {
	existFileAllowPolicy := false
	existNetworkAllowPolicy := false
	existCapabilitiesAllowPolicy := false

	if log.Result == "Passed" || log.Result == "Operation not permitted" || log.Result == "Permission denied" {
		fd.SecurityPoliciesLock.RLock()

		key := cfg.GlobalCfg.Host

		if log.NamespaceName != "" && log.PodName != "" {
			key = log.NamespaceName + "_" + log.PodName
		}

		secPolicies := fd.SecurityPolicies[key].Policies
		for _, secPolicy := range secPolicies {
			if secPolicy.Action == "Allow" || secPolicy.Action == "Audit (Allow)" {
				if secPolicy.Operation == "Process" || secPolicy.Operation == "File" {
					existFileAllowPolicy = true
				} else if secPolicy.Operation == "Network" {
					existNetworkAllowPolicy = true
				} else if secPolicy.Operation == "Capabilities" {
					existCapabilitiesAllowPolicy = true
				}

				if fd.DefaultPostures[log.NamespaceName].FileAction == "allow" {
					continue
				}
			}

			firstLogResource := strings.Split(log.Resource, " ")[0]
			firstLogResourceDir := getDirectoryPart(firstLogResource)
			firstLogResourceDirCount := strings.Count(firstLogResourceDir, "/")

			switch log.Operation {
			case "Process", "File":
				if secPolicy.Operation != log.Operation {
					continue
				}

				// match sources
				if (!secPolicy.IsFromSource) || (secPolicy.IsFromSource && (secPolicy.Source == log.ParentProcessName || secPolicy.Source == log.ProcessName)) {
					matchedRegex := false

					switch secPolicy.ResourceType {
					case "Glob":
						// Match using a globbing syntax very similar to the AppArmor's
						matchedRegex, _ = filepath.Match(secPolicy.Resource, log.Resource) // pattern (secPolicy.Resource) -> string (log.Resource)
					case "Regexp":
						if secPolicy.Regexp != nil {
							// Match using compiled regular expression
							matchedRegex = secPolicy.Regexp.MatchString(log.Resource) // regexp (secPolicy.Regexp) -> string (log.Resource)
						}
					}

					// match resources
					if matchedRegex || (secPolicy.ResourceType == "Path" && secPolicy.Resource == firstLogResource) ||
						(secPolicy.ResourceType == "Directory" && strings.HasPrefix(firstLogResourceDir, secPolicy.Resource) &&
							((!secPolicy.Recursive && firstLogResourceDirCount == strings.Count(secPolicy.Resource, "/")) ||
								(secPolicy.Recursive && firstLogResourceDirCount >= strings.Count(secPolicy.Resource, "/")))) {

						matchedFlags := false

						if secPolicy.ReadOnly && log.Resource != "" && secPolicy.OwnerOnly && log.MergedDir != "" {
							// read only && owner only
							if strings.Contains(log.Data, "O_RDONLY") && strconv.Itoa(int(log.UID)) == getFileProcessUID(log.MergedDir+log.Resource) {
								matchedFlags = true
							}
						} else if secPolicy.ReadOnly && log.Resource != "" {
							// read only
							if strings.Contains(log.Data, "O_RDONLY") {
								matchedFlags = true
							}
						} else if secPolicy.OwnerOnly && log.MergedDir != "" {
							// owner only
							if strconv.Itoa(int(log.UID)) == getFileProcessUID(log.MergedDir+log.Resource) {
								matchedFlags = true
							}
						} else {
							// ! read only && ! owner only
							matchedFlags = true
						}

						if matchedFlags && (secPolicy.Action == "Allow" || secPolicy.Action == "Audit (Allow)") && log.Result == "Passed" {
							// allow policy or allow policy with audit mode
							// matched source + matched resource + matched flags + matched action + expected result -> going to be skipped

							log.Type = "MatchedPolicy"

							log.PolicyName = secPolicy.PolicyName
							log.Severity = secPolicy.Severity

							if len(secPolicy.Tags) > 0 {
								log.Tags = strings.Join(secPolicy.Tags[:], ",")
								log.ATags = secPolicy.Tags
							}

							if len(secPolicy.Message) > 0 {
								log.Message = secPolicy.Message
							}

							if log.PolicyEnabled == tp.KubeArmorPolicyAudited {
								log.Enforcer = "eBPF Monitor"
							} else {
								log.Enforcer = fd.Enforcer
							}

							log.Action = "Allow"

							continue
						}

						if matchedFlags && secPolicy.Action == "Audit" && log.Result == "Passed" {
							// audit policy
							// matched source + matched resource + matched flags + matched action + expected result -> alert (audit log)

							log.Type = "MatchedPolicy"

							log.PolicyName = secPolicy.PolicyName
							log.Severity = secPolicy.Severity

							if len(secPolicy.Tags) > 0 {
								log.Tags = strings.Join(secPolicy.Tags[:], ",")
								log.ATags = secPolicy.Tags
							}

							if len(secPolicy.Message) > 0 {
								log.Message = secPolicy.Message
							}

							log.Enforcer = "eBPF Monitor"
							log.Action = secPolicy.Action

							continue
						}

						if (secPolicy.Action == "Block" && log.Result != "Passed") ||
							(matchedFlags && (!secPolicy.OwnerOnly && !secPolicy.ReadOnly) && secPolicy.Action == "Audit (Block)" && log.Result == "Passed") ||
							(!matchedFlags && (secPolicy.OwnerOnly || secPolicy.ReadOnly) && secPolicy.Action == "Audit (Block)" && log.Result == "Passed") {
							// block policy or block policy with audit mode
							// matched source + matched resource + matched action + expected result -> alert

							log.Type = "MatchedPolicy"

							log.PolicyName = secPolicy.PolicyName
							log.Severity = secPolicy.Severity

							if len(secPolicy.Tags) > 0 {
								log.Tags = strings.Join(secPolicy.Tags[:], ",")
								log.ATags = secPolicy.Tags
							}

							if len(secPolicy.Message) > 0 {
								log.Message = secPolicy.Message
							}

							if log.PolicyEnabled == tp.KubeArmorPolicyAudited {
								log.Enforcer = "eBPF Monitor"
							} else {
								log.Enforcer = fd.Enforcer
							}

							log.Action = secPolicy.Action

							continue
						}

						if matchedFlags && secPolicy.Action == "Allow" && log.Result != "Passed" {
							// It's possible there are additional rules in the Security Policy resulting in the block else we deem it as default posture anyway
							continue
						}
					}

					if secPolicy.Action == "Allow" && log.Result != "Passed" {
						// matched source + !(matched resource) + action = allow + result = blocked -> default posture / allow policy violation

						log.Type = "MatchedPolicy"

						log.PolicyName = "DefaultPosture"

						log.Severity = ""
						log.Tags = ""
						log.ATags = []string{}
						log.Message = ""

						log.Enforcer = "eBPF Monitor"
						log.Action = "Block"

						continue
					}

					if secPolicy.Action == "Audit (Allow)" && log.Result == "Passed" {
						// matched source + !(matched resource) + action = audit (allow) + result = passed -> default posture / allow policy violation (audit mode)

						log.Type = "MatchedPolicy"

						log.PolicyName = "DefaultPosture"

						log.Severity = ""
						log.Tags = ""
						log.ATags = []string{}
						log.Message = ""

						log.Enforcer = "eBPF Monitor"

						if fd.DefaultPostures[log.NamespaceName].FileAction == "block" {
							log.Action = "Audit (Block)"
						} else { // fd.DefaultPostures[log.NamespaceName].FileAction == "audit"
							log.Action = "Audit"
						}

						continue
					}
				}

				if fd.DefaultPostures[log.NamespaceName].FileAction == "block" && secPolicy.Action == "Audit (Allow)" && log.Result == "Passed" {
					// defaultPosture = block + audit mode

					log.Type = "MatchedPolicy"

					log.PolicyName = "DefaultPosture"

					log.Severity = ""
					log.Tags = ""
					log.ATags = []string{}
					log.Message = ""

					log.Enforcer = "eBPF Monitor"
					log.Action = "Audit (Block)"
				}

				if fd.DefaultPostures[log.NamespaceName].FileAction == "audit" && (secPolicy.Action == "Allow" || secPolicy.Action == "Audit (Allow)") && log.Result == "Passed" {
					// defaultPosture = audit

					log.Type = "MatchedPolicy"

					log.PolicyName = "DefaultPosture"

					log.Severity = ""
					log.Tags = ""
					log.ATags = []string{}
					log.Message = ""

					log.Enforcer = "eBPF Monitor"
					log.Action = "Audit"
				}

			case "Network":
				if secPolicy.Operation != log.Operation {
					continue
				}

				// match sources
				if (!secPolicy.IsFromSource) || (secPolicy.IsFromSource && (secPolicy.Source == log.ParentProcessName || secPolicy.Source == log.ProcessName)) {
					skip := false

					for _, matchProtocol := range strings.Split(secPolicy.Resource, ",") {
						if skip {
							break
						}

						// match resources
						if strings.Contains(log.Resource, matchProtocol) {
							if (secPolicy.Action == "Allow" || secPolicy.Action == "Audit (Allow)") && log.Result == "Passed" {
								// allow policy or allow policy with audit mode
								// matched source + matched resource + matched action + expected result -> going to be skipped

								log.Type = "MatchedPolicy"

								log.PolicyName = secPolicy.PolicyName
								log.Severity = secPolicy.Severity

								if len(secPolicy.Tags) > 0 {
									log.Tags = strings.Join(secPolicy.Tags[:], ",")
									log.ATags = secPolicy.Tags
								}

								if len(secPolicy.Message) > 0 {
									log.Message = secPolicy.Message
								}

								if log.PolicyEnabled == tp.KubeArmorPolicyAudited {
									log.Enforcer = "eBPF Monitor"
								} else {
									log.Enforcer = fd.Enforcer
								}

								log.Action = "Allow"

								skip = true
								continue
							}

							if secPolicy.Action == "Audit" && log.Result == "Passed" {
								// audit policy
								// matched source + matched resource + matched action + expected result -> alert (audit log)

								log.Type = "MatchedPolicy"

								log.PolicyName = secPolicy.PolicyName
								log.Severity = secPolicy.Severity

								if len(secPolicy.Tags) > 0 {
									log.Tags = strings.Join(secPolicy.Tags[:], ",")
									log.ATags = secPolicy.Tags
								}

								if len(secPolicy.Message) > 0 {
									log.Message = secPolicy.Message
								}

								log.Enforcer = "eBPF Monitor"
								log.Action = secPolicy.Action

								skip = true
								continue
							}

							if (secPolicy.Action == "Block" && log.Result != "Passed") ||
								(secPolicy.Action == "Audit (Block)" && log.Result == "Passed") {
								// block policy or block policy with audit mode
								// matched source + matched resource + matched action + expected result -> alert

								log.Type = "MatchedPolicy"

								log.PolicyName = secPolicy.PolicyName
								log.Severity = secPolicy.Severity

								if len(secPolicy.Tags) > 0 {
									log.Tags = strings.Join(secPolicy.Tags[:], ",")
								}

								if len(secPolicy.Message) > 0 {
									log.Message = secPolicy.Message
								}

								if log.PolicyEnabled == tp.KubeArmorPolicyAudited {
									log.Enforcer = "eBPF Monitor"
								} else {
									log.Enforcer = fd.Enforcer
								}

								log.Action = secPolicy.Action

								skip = true
								continue
							}
						}
					}

					if skip {
						continue
					}

					if secPolicy.Action == "Allow" && log.Result != "Passed" {
						// matched source + !(matched resource) + action = allow + result = blocked -> allow policy violation

						log.Type = "MatchedPolicy"

						log.PolicyName = "DefaultPosture"

						log.Severity = ""
						log.Tags = ""
						log.Message = ""

						log.Enforcer = "eBPF Monitor"
						log.Action = "Block"

						continue
					}

					if secPolicy.Action == "Audit (Allow)" && log.Result == "Passed" {
						// matched source + !(matched resource) + action = audit (allow) + result = passed -> allow policy violation (audit mode)

						log.Type = "MatchedPolicy"

						log.PolicyName = "DefaultPosture"

						log.Severity = ""
						log.Tags = ""
						log.Message = ""

						log.Enforcer = "eBPF Monitor"

						if fd.DefaultPostures[log.NamespaceName].NetworkAction == "block" {
							log.Action = "Audit (Block)"
						} else { // fd.DefaultPostures[log.NamespaceName].NetworkAction == "audit"
							log.Action = "Audit"
						}

						continue
					}
				}

				if fd.DefaultPostures[log.NamespaceName].NetworkAction == "block" && secPolicy.Action == "Audit (Allow)" && log.Result == "Passed" {
					// defaultPosture = block + audit mode

					log.Type = "MatchedPolicy"

					log.PolicyName = "DefaultPosture"

					log.Severity = ""
					log.Tags = ""
					log.Message = ""

					log.Enforcer = "eBPF Monitor"
					log.Action = "Audit (Block)"
				}

				if fd.DefaultPostures[log.NamespaceName].NetworkAction == "audit" && (secPolicy.Action == "Allow" || secPolicy.Action == "Audit (Allow)") && log.Result == "Passed" {
					// defaultPosture = audit

					log.Type = "MatchedPolicy"

					log.PolicyName = "DefaultPosture"

					log.Severity = ""
					log.Tags = ""
					log.Message = ""

					log.Enforcer = "eBPF Monitor"
					log.Action = "Audit"
				}
			case "Syscall":
				if secPolicy.Operation != log.Operation {
					continue
				}
				//Get syscall
				syscallName := strings.Split(strings.Split(log.Data, " ")[0], "SYS_")[1]
				//Get syscall Source
				syscallSource := strings.Split(log.Source, " ")[0]
				matchedRule := false
				if syscallName == secPolicy.ResourceType {
					matchPath := false
					fromSource := false
					if secPolicy.IsFromSource &&
						(((strings.HasPrefix(syscallSource, secPolicy.Source) && secPolicy.Source[len(secPolicy.Source)-1] == '/') && // match dir
							(secPolicy.Recursive || !strings.Contains(syscallSource[len(secPolicy.Source):], "/"))) || // handle recursive dir
							secPolicy.Source == syscallSource) { // match file
						fromSource = true
					}

					if len(secPolicy.Resource) > 0 &&
						((secPolicy.Resource[len(secPolicy.Resource)-1] == '/' && ((strings.HasPrefix(log.Resource, secPolicy.Resource) && secPolicy.ReadOnly) || secPolicy.Resource[:len(secPolicy.Resource)-1] == log.Resource)) || //match dir
							secPolicy.Resource == log.Resource) { // match path
						matchPath = true
					}
					matchedRule = (len(secPolicy.Resource) == 0 || matchPath) && (!secPolicy.IsFromSource || fromSource)

					if matchedRule {
						log.Type = "MatchedPolicy"
						log.PolicyName = secPolicy.PolicyName
						log.Severity = secPolicy.Severity
						if len(secPolicy.Tags) > 0 {
							log.Tags = strings.Join(secPolicy.Tags[:], ",")
						}

						if len(secPolicy.Message) > 0 {
							log.Message = secPolicy.Message
						}
					}
				}

			}
		}

		fd.SecurityPoliciesLock.RUnlock()

		if log.PolicyName == "" && log.Result != "Passed" {
			// default posture (block) or native policy
			// no matched policy, but result = blocked -> default posture

			log.Type = "MatchedPolicy"

			log.PolicyName = "DefaultPosture"

			log.Severity = ""
			log.Tags = ""
			log.ATags = []string{}
			log.Message = ""

			log.Enforcer = fd.Enforcer
			log.Action = "Block"
		}
	}

	if log.ContainerID != "" { // container
		if log.Type == "" {
			// defaultPosture (audit) or container log

			fd.DefaultPosturesLock.Lock()

			if _, ok := fd.DefaultPostures[log.NamespaceName]; !ok {
				globalDefaultPosture := tp.DefaultPosture{
					FileAction:         cfg.GlobalCfg.DefaultFilePosture,
					NetworkAction:      cfg.GlobalCfg.DefaultNetworkPosture,
					CapabilitiesAction: cfg.GlobalCfg.DefaultCapabilitiesPosture,
				}
				fd.DefaultPostures[log.NamespaceName] = globalDefaultPosture
			}

			fd.DefaultPosturesLock.Unlock()

			if log.Operation == "Process" {
				if setLogFields(&log, existFileAllowPolicy, fd.DefaultPostures[log.NamespaceName].FileAction, log.ProcessVisibilityEnabled, true) {
					return log
				}
			} else if log.Operation == "File" {
				if setLogFields(&log, existFileAllowPolicy, fd.DefaultPostures[log.NamespaceName].FileAction, log.FileVisibilityEnabled, true) {
					return log
				}
			} else if log.Operation == "Network" {
				if setLogFields(&log, existNetworkAllowPolicy, fd.DefaultPostures[log.NamespaceName].NetworkAction, log.NetworkVisibilityEnabled, true) {
					return log
				}
			} else if log.Operation == "Capabilities" {
				if setLogFields(&log, existCapabilitiesAllowPolicy, fd.DefaultPostures[log.NamespaceName].CapabilitiesAction, log.CapabilitiesVisibilityEnabled, true) {
					return log
				}
			} else if log.Operation == "Syscall" {
				if setLogFields(&log, false, "", true, true) {
					return log
				}
			}

		} else if log.Type == "MatchedPolicy" {
			if log.Action == "Allow" && log.Result == "Passed" {
				return tp.Log{}
			}

			return log
		}
	} else { // host
		if log.Type == "" {
			// host log

			if log.Operation == "Process" {
				if setLogFields(&log, existFileAllowPolicy, "allow", fd.Node.ProcessVisibilityEnabled, false) {
					return log
				}
			} else if log.Operation == "File" {
				if setLogFields(&log, existFileAllowPolicy, "allow", fd.Node.FileVisibilityEnabled, false) {
					return log
				}
			} else if log.Operation == "Network" {
				if setLogFields(&log, existNetworkAllowPolicy, "allow", fd.Node.NetworkVisibilityEnabled, false) {
					return log
				}
			} else if log.Operation == "Capabilities" {
				if setLogFields(&log, existCapabilitiesAllowPolicy, "allow", fd.Node.CapabilitiesVisibilityEnabled, false) {
					return log
				}
			}

		} else if log.Type == "MatchedPolicy" {
			log.Type = "MatchedHostPolicy"

			if log.Action == "Allow" && log.Result == "Passed" {
				return tp.Log{}
			}

			return log
		}
	}

	return tp.Log{}
}
