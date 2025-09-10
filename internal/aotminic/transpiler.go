package aotminic

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

func TranspileToC(src string, name string) (string, error) {
	var b bytes.Buffer

	// C prologue: standard headers + runtime prototypes
	b.WriteString("#include <stdio.h>\n#include <stdint.h>\n#include <stdlib.h>\n\n")
	b.WriteString("extern int __argc; extern char** __argv;\n")
	b.WriteString("long print(long x);\n")
	b.WriteString("long argi(long idx);\n")
	b.WriteString("long time_ns(void);\n")
	b.WriteString("long print_time_ns(long ns);\n\n")

	sc := bufio.NewScanner(strings.NewReader(src))
	spaceRE := regexp.MustCompile(`[ \t]+`)
	trueRE := regexp.MustCompile(`\btrue\b`)
	falseRE := regexp.MustCompile(`\bfalse\b`)

	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimRight(line, "\r")
		line = spaceRE.ReplaceAllString(line, " ")
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			b.WriteString("/* " + strings.TrimSpace(line)[2:] + " */\n")
			continue
		}
		trim := strings.TrimSpace(line)
		trim = trueRE.ReplaceAllString(trim, "1")
		trim = falseRE.ReplaceAllString(trim, "0")

		switch {
		case strings.HasPrefix(trim, "fn "):
			header := strings.TrimPrefix(trim, "fn ")
			header = strings.TrimSpace(header)
			if strings.HasPrefix(header, "main(") {
				b.WriteString("int main(int argc, char** argv) { __argc=argc; __argv=argv;\n")
				continue
			}
			fnSig := header
			namePart, paramsPart, tail := splitSignature(fnSig)
			paramsTyped := typeParamsAsLong(paramsPart)
			trim = "long " + namePart + "(" + paramsTyped + ") " + tail
			b.WriteString(trim + "\n")
			continue

		case strings.HasPrefix(trim, "let "):
			rest := strings.TrimPrefix(trim, "let ")
			trim = "long " + rest
			if !strings.HasSuffix(trim, ";") && strings.HasSuffix(trim, ")") {
				trim += ";"
			}
			b.WriteString(trim + "\n")
			continue

		case strings.HasPrefix(trim, "return "):
			b.WriteString(trim + "\n")
			continue

		case strings.HasPrefix(trim, "if ") || strings.HasPrefix(trim, "if(") ||
			strings.HasPrefix(trim, "while ") || strings.HasPrefix(trim, "while(") ||
			trim == "}" || trim == "else {" || strings.HasPrefix(trim, "else ") || trim == "{":
			b.WriteString(trim + "\n")
			continue

		default:
			if len(strings.TrimSpace(trim)) == 0 {
				b.WriteString("\n")
				continue
			}
			if !strings.HasSuffix(trim, "}") && !strings.HasSuffix(trim, "{") && !strings.HasSuffix(trim, ";") {
				trim += ";"
			}
			b.WriteString(trim + "\n")
		}
	}
	return b.String(), nil
}

func splitSignature(sig string) (name string, params string, tail string) {
	sig = strings.TrimSpace(sig)
	open := strings.Index(sig, "(")
	close := strings.LastIndex(sig, ")")
	if open == -1 || close == -1 || close < open {
		parts := strings.SplitN(sig, " ", 2)
		if len(parts) == 2 {
			return parts[0], "", parts[1]
		}
		return sig, "", "{"
	}
	name = strings.TrimSpace(sig[:open])
	params = strings.TrimSpace(sig[open+1 : close])

	rest := strings.TrimSpace(sig[close+1:])
	if rest == "" {
		rest = "{"
	}
	return name, params, rest
}

func typeParamsAsLong(params string) string {
	params = strings.TrimSpace(params)
	if params == "" {
		return ""
	}
	parts := strings.Split(params, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.Contains(p, " ") {
			out = append(out, p)
		} else {
			out = append(out, "long "+p)
		}
	}
	return strings.Join(out, ", ")
}