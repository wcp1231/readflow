import React from 'react'

import { Article } from '../../models'
import OfflineLink from './OfflineLink'
import ShareLink from './ShareLink'
import DownloadAsLink from './DownloadAsLink'
import OutgoingWebhooksMenuItems from './OutgoingWebhooksMenuItems'
import DrawerMenu from '../../../components/DrawerMenu'

interface Props {
  article: Article
  keyboard?: boolean
}

export default (props: Props) => {
  const nvg: any = window.navigator
  const title = 'More actions...'

  return (
    <DrawerMenu title={title}>
      <ul>
        {nvg.share && (
          <li>
            <ShareLink {...props} />
          </li>
        )}
        <li>
          <DownloadAsLink {...props} />
        </li>
        <li>
          <OfflineLink {...props} />
        </li>
        <OutgoingWebhooksMenuItems {...props} />
      </ul>
    </DrawerMenu>
  )
}
