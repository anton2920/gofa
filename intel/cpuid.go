package intel

//go:nosplit
func CPUID(eax, ecx uint32) (reax, rebx, recx, redx uint32)
