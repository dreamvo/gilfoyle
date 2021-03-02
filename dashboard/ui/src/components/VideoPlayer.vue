<template>
  <div>
    <div v-show="!loading">
      <video ref="player" class="video-js" :width="width" :height="height">
        <p class="vjs-no-js">
          To view this video please enable JavaScript, and consider upgrading to
          a web browser that
          <a href="https://videojs.com/html5-video-support/" target="_blank"
            >supports HTML5 video</a
          >
        </p>
      </video>
    </div>
    <div v-show="loading">Loading media player...</div>
  </div>
</template>

<style>
.video-js {
  width: 100%;
}
</style>

<script lang="ts">
import Vue from "vue";
import "video.js/dist/video-js.min.css";
import videojs from "video.js";
import "videojs-hls-quality-selector";
import "videojs-contrib-quality-levels";

interface Data {
  loading: boolean;
  player: videojs.Player;
}

export default Vue.extend({
  name: "VideoPlayer",
  props: {
    width: {
      type: String,
      default: "100%"
    },
    height: {
      type: String,
      default: "360"
    },
    sources: {
      type: Array,
      default: () => []
    },
    posterURL: {
      type: String,
      default: ""
    }
  },
  watch: {
    sources(/*sources: Source[]*/) {
      this.stream();
    }
  },
  data: (): Data => ({
    loading: true,
    player: videojs.getPlayer("video")
  }),
  methods: {
    stream() {
      const video = this.$refs.player as HTMLMediaElement;

      this.player = videojs(
        video,
        {
          aspectRatio: "16:9",
          autoplay: false,
          loop: false,
          controls: true,
          // poster: '/app/img/default_media.0f638ccd.jpeg',
          bigPlayButton: true,
          sources: this.sources as videojs.Tech.SourceObject[],
          preload: "auto",
          playbackRates: [0.5, 1, 1.5, 2],
          plugins: {
            hlsQualitySelector: {
              displayCurrentQuality: true
            }
          }
        },
        () => {
          console.info("Player is ready");
        }
      );
    }
  },
  async created() {
    await new Promise(resolve => setTimeout(resolve, 100));

    this.stream();

    this.loading = false;
  }
});
</script>
