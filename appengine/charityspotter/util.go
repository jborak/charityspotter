package charityspotter

import (
  "bytes"
)

type ClosingBuffer struct { 
  *bytes.Buffer
} 

func (cb *ClosingBuffer) Close() (error) { 
  return nil
}
