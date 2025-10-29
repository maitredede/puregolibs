package libevdi

const (
	libEvdiVersionMajor = 1
	libEvdiVersionMinor = 14
	libEvdiVersionPatch = 1

	evdiModuleCompatibilityVersionMajor = 1
	evdiModuleCompatibilityVersionMinor = 9
	evdiModuleCompatibilityVersionPatch = 0
)

type LibVersion struct {
	Major      int
	Minor      int
	PatchLevel int
}

func GetLibVersion() LibVersion {
	return LibVersion{
		Major:      libEvdiVersionMajor,
		Minor:      libEvdiVersionMinor,
		PatchLevel: libEvdiVersionPatch,
	}
}
