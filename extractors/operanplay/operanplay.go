package operanplay

import (
//   "encoding/json"
 //  "strings"
  // "fmt"
 "sort"

   "github.com/pkg/errors"
   "github.com/iawia002/lia/array"

   "github.com/fabe-xx/lux/extractors"
   "github.com/fabe-xx/lux/extractors/vimeo"
   "github.com/fabe-xx/lux/request"
   "github.com/fabe-xx/lux/utils"
)

func init() {
   extractors.Register("operanplay", New())
}

type extractor struct{}

func New() extractors.Extractor {
   return extractor{}
}

// Extract is the main function to extract the data.
func (extractor) Extract(url string, option extractors.Options) (ret []*extractors.Data, err error) {
   var (
      html string
   )
   
      html, err = request.Get(url, url, nil)
      if err != nil {
         return nil, errors.WithStack(err)
      }
   
   matches := utils.MatchAll(html, `data-vimeo-url\s*=\s*"(https://vimeo.com/\d+)`)
   if len(matches) == 0 {
      return nil, errors.WithStack(extractors.ErrURLParseFailed)
   }
   
   m := map[string]struct{}{}
   vimeos := []string {}
   
   for _, vim := range matches {
      if _, ok := m[vim[1]]; !ok {
         m[vim[1]] = struct{}{}
         vimeos = append(vimeos,vim[1])   
      }
   }
   
   sort.Strings( vimeos )
   
   vimExt := vimeo.New()
   
   needDownloadItems := utils.NeedDownloadList(option.Items, option.ItemStart, option.ItemEnd, len(vimeos))
   for i, vim := range vimeos {
      if !array.ItemInArray(i+1, needDownloadItems) {
         continue
      }
      if bits, err := vimExt.Extract(vim,option); err != nil {
         return nil, errors.WithStack(err)
      } else {
         for _, bit := range bits {
            bit.Site = "Operan Play operanplay.se"  
            bit.URL = url
    //        m := map[string]*extractors.Stream {}
      //      for n, d := range bit.Streams {
        //       m[fmt.Sprint(  i+1,"/",n)] = d
          //  }
            //bit.Streams = m
            ret = append( ret, bit )
         }
      }
   }

   return
}


/*
   for _, video := range vimeoData.Request.Files.Progressive {
      size, err = request.Size(video.URL, url)
      if err != nil {
         return nil, errors.WithStack(err)
      }
      urlData := &extractors.Part{
         URL:  video.URL,
         Size: size,
         Ext:  "mp4",
      }
      streams[video.Profile] = &extractors.Stream{
         Parts:   []*extractors.Part{urlData},
         Size:    size,
         Quality: video.Quality,
      }
   }

   return []*extractors.Data{
      {
         Site:    "Vimeo vimeo.com",
         Title:   vimeoData.Video.Title,
         Type:    extractors.DataTypeVideo,
         Streams: streams,
         URL:     url,
      },
   }, nil

*/