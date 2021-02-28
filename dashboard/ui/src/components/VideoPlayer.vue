<template>
  <video
    ref="player"
    id="my-video"
    class="video-js"
    controls
    :muted="false"
    preload="auto"
    :width="width"
    :height="height"
    :poster="posterURL"
    data-setup="{}"
  >
    <source
      v-for="(source, index) of sources"
      :key="index"
      :src="source.src"
      :type="source.type"
    />
    <p class="vjs-no-js">
      To view this video please enable JavaScript, and consider upgrading to a
      web browser that
      <a href="https://videojs.com/html5-video-support/" target="_blank"
        >supports HTML5 video</a
      >
    </p>
  </video>
</template>

<style>
.video-js {
  max-width: 100%;
}
</style>

<script lang="ts">
import Vue from "vue";
import Hls from "hls.js";

interface Source {
  src: string;
  type: string;
}

interface Data {}

export default Vue.extend({
  name: "VideoPlayer",
  props: {
    width: {
      type: String,
      default: "100%"
    },
    height: {
      type: String,
      default: "auto"
    },
    sources: {
      type: Array,
      default: [] as Source[]
    },
    posterURL: {
      type: String,
      default: ""
    }
  },
  watch: {
    sources(sources: Source[]) {
      this.stream(sources[0].src);
    }
  },
  data: (): Data => ({}),
  methods: {
    stream(source: string) {
      const hls = new Hls();
      const video = this.$refs.player as HTMLMediaElement;

      hls.loadSource(source);
      hls.attachMedia(video);
    }
  },
  created() {
    if (this.sources.length > 0) {
      this.stream(this.sources[0].src);
    }
  }
});
</script>
