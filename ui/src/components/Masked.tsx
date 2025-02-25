import React, { createRef } from 'react'
import { ButtonIcon } from '.'

interface Props {
  value: string
}

export const Masked = ({ value }: Props) => {
  const ref = createRef<HTMLInputElement>()

  function copyToClipboard() {
    if (ref.current) {
      ref.current.select()
      document.execCommand('copy')
    }
  }

  return (
    <div>
      <input type="text" readOnly ref={ref} className="masked" value={value} />
      <ButtonIcon icon="file_copy" onClick={copyToClipboard} title="Copy to the clipboard" />
    </div>
  )
}
