import React, { FC, ImgHTMLAttributes } from 'react'

import { API_BASE_URL } from '../../constants'

const proxifyImageURL = (url: string, width: number) =>
  `${API_BASE_URL}/img?url=${encodeURIComponent(url)}&width=${width}`

export const ArticleImage: FC<ImgHTMLAttributes<HTMLImageElement>> = ({ src, ...attrs }) => {
  if (src && src.match(/^https?:\/\//)) {
    attrs.srcSet = `${proxifyImageURL(src, 320)} 320w, ${proxifyImageURL(src, 767)} 767w`
  }
  return (
    <img
      {...attrs}
      sizes="(max-width: 767px) 767px, 320px"
      src={src}
      onError={(e) => (e.currentTarget.style.display = 'none')}
    />
  )
}
