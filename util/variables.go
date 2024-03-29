package util

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/eris-ltd/eris-pm/definitions"
)

func PreProcess(toProcess string, do *definitions.Do) (string, error) {
	// $block.... $account.... etc. should be caught. hell$$o should not
	catchEr := regexp.MustCompile("^\\$(.*)$")

	// If there's a match then run through the replacement process
	if catchEr.MatchString(toProcess) {
		// find what we need to catch.
		jobName := catchEr.FindStringSubmatch(toProcess)[1]

		// first parse the reserved words.
		if strings.Contains(jobName, "block") {
			return replaceBlockVariable(jobName, do)
		}

		// second we loop through the jobNames to do a result replace
		for _, job := range do.Package.Jobs {
			if string(jobName) == job.JobName {
				logger.Debugf("Fixing Variables =>\t\t$%s:%s\n", string(jobName), job.JobResult)
				return job.JobResult, nil
			}
		}
	}

	// if no matches, return original
	return toProcess, nil
}

func replaceBlockVariable(toReplace string, do *definitions.Do) (string, error) {
	block, err := ChainStatus("latest_block_height", do)
	if err != nil {
		return "", err
	}

	if toReplace == "block" {
		return block, nil
	}

	catchEr := regexp.MustCompile("block\\+(\\d*)")
	if catchEr.MatchString(toReplace) {
		height := catchEr.FindStringSubmatch(toReplace)[1]
		h1, err := strconv.Atoi(height)
		if err != nil {
			return "", err
		}
		h2, err := strconv.Atoi(block)
		if err != nil {
			return "", err
		}
		height = strconv.Itoa(h1 + h2)
		return height, nil
	}

	catchEr = regexp.MustCompile("block\\-(\\d*)")
	if catchEr.MatchString(toReplace) {
		height := catchEr.FindStringSubmatch(toReplace)[1]
		h1, err := strconv.Atoi(height)
		if err != nil {
			return "", err
		}
		h2, err := strconv.Atoi(block)
		if err != nil {
			return "", err
		}
		height = strconv.Itoa(h1 - h2)
		return height, nil
	}

	return toReplace, nil
}
