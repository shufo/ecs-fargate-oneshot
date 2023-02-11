package tasks

func ValidateCombinationOfCpuAndMemory(cpu uint64, memory uint64) bool {
	if cpu == 0 && memory == 0 {
		return true
	}

	if cpu == 256 {
		if memory >= 512 && memory <= 2048 {
			return true
		}

		return false
	}

	if cpu == 512 {
		if memory >= 1024 && memory <= 4096 {
			return true
		}

		return false
	}

	if cpu == 1024 {
		if memory >= 2048 && memory <= 8192 {
			return true
		}

		return false
	}

	if cpu == 2048 {
		if memory >= 4096 && memory <= 16384 {
			return true
		}

		return false
	}

	if cpu == 4096 {
		if memory >= 8192 && memory <= 30720 {
			return true
		}

		return false
	}

	if cpu == 8192 {
		if memory >= 16384 && memory <= 61440 {
			return true
		}

		return false
	}

	if cpu == 16384 {
		if memory >= 32768 && memory <= 122880 {
			return true
		}

		return false
	}

	return false
}
