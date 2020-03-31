package metapod

import (
	"github.com/RainwayApp/metapod/windows"
)

// Create adds the payload to the target executeable and returns the result
func Create(peFile []byte, payload []byte) ([]byte, error) {
	portableExecutable, err := windows.GetPortableExecutable(peFile)

	if err != nil {
		return []byte{}, err
	}
	var targetExecutable = windows.TargetExecutable{*portableExecutable}
	contents, err := targetExecutable.CreateFromTemplate(payload)
	if err != nil {
		return []byte{}, err
	}

	return contents, nil
}

// Open gets the payload from a file
// rawPayload may return nil with no error - this means that the payload did
// not exist
func Open(peFile []byte) ([]byte, error) {
	portableExecutable, err := windows.GetPortableExecutable(peFile)

	if err != nil {
		return []byte{}, err
	}

	var targetExecutable = windows.TargetExecutable{*portableExecutable}
	_, rawPayload, err := targetExecutable.GetPayload()

	if err != nil {
		return []byte{}, err
	}

	return rawPayload, nil
}
