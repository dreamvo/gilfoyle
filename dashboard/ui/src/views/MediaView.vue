<template>
  <v-row>
    <v-col md="6" sm="12">
      <v-card flat outlined :loading="loading">
        <v-card-title>{{ media.title }}</v-card-title>
        <v-card-subtitle class="mt-1">
          <div><strong>UUID</strong> : {{ media.id }}</div>
          <div><strong>Status</strong> :
            <v-chip :color="statusLevel" x-small>{{ media.status }}</v-chip>
          </div>
          <div>Created on {{ media.created_at }}</div>
        </v-card-subtitle>

        <v-card-actions>
          <v-btn
              v-if="media.status === 'AwaitingUpload'"
              text
              color="primary"
              small
          >
            <v-icon size="16" dark class="mr-1">
              mdi-upload
            </v-icon>
            Upload media file
          </v-btn
          >

          <v-btn small text color="secondary" dark class="ml-2" :to="{name:'MediaUpdate', id: media.id}">
            <v-icon size="16" dark class="mr-1">
              mdi-pencil
            </v-icon>
            Update
          </v-btn>

          <DeleteModal btn-small btn-text="Delete" :action="deleteMedia"/>
        </v-card-actions>
      </v-card>

      <v-card flat outlined :loading="loading" class="mt-6" v-if="media.edges && media.edges.media_files">
        <v-card-title>Renditions</v-card-title>

        <v-simple-table>
          <template v-slot:default>
            <thead>
            <tr>
              <th class="text-left">
                Name
              </th>
              <th class="text-left">
                Format
              </th>
              <th class="text-left">
                Resolution
              </th>
              <th class="text-left">
                Target bandwidth
              </th>
              <th class="text-left">
                Framerate
              </th>
              <th class="text-left">
                Duration
              </th>
              <th class="text-left">
                Media type
              </th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="item in media.edges.media_files" :key="item.id">
              <td>{{ item.rendition_name }}</td>
              <td>{{ item.format }}</td>
              <td>
                {{ item.resolution_width }}x{{ item.resolution_height }}
              </td>
              <td>{{ item.target_bandwidth }}</td>
              <td>{{ item.framerate }}</td>
              <td>{{ item.duration_seconds }}</td>
              <td>{{ item.media_type }}</td>
            </tr>
            </tbody>
          </template>
        </v-simple-table>
      </v-card>
    </v-col>

    <v-col md="6" sm="12">
      <v-card flat outlined :loading="loading">
        <v-card-title>Streaming</v-card-title>

        <v-card-text v-if="!streamReady"
        >This media is not yet available for streaming.
        </v-card-text>
        <v-card-text v-show="streamReady">
          <video ref="player" controls muted loop width="100%"></video>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import {Vue} from "vue-property-decorator";
import {AxiosResponse} from "axios";
import {DataResponse, Media} from "../types";
import axios from "../services/axios";
import Hls from "hls.js";
import DeleteModal from "../components/DeleteModal.vue";

interface Data {
  loading: boolean;
  streamReady: boolean;
  media: Media;
}

export default Vue.extend({
  components: {DeleteModal},
  data: () => ({
    loading: true,
    streamReady: false,
    media: {} as Media
  }),
  methods: {
    stream() {
      const hls = new Hls();
      const stream = `${axios.defaults.baseURL}/medias/${this.media.id}/stream/master.m3u8`;
      const video: any = this.$refs.player;

      hls.loadSource(stream);
      hls.attachMedia(video);

      hls.on(Hls.Events.MANIFEST_PARSED, function () {
        video.play();
      });
    },
    async deleteMedia() {
      await axios.delete(
          `/medias/${this.$route.params.id}`
      );

      await this.$router.push({name: 'MediaAll'})
    }
  },
  computed: {
    statusLevel(): string {
      const status: { [key: string]: string } = {
        'Ready': 'success',
        'Errored': 'error',
        'AwaitingUpload': 'warning',
      }

      return status[this.media.status] || ''
    }
  },
  async created() {
    const res: AxiosResponse<DataResponse<Media>> = await axios.get(
        `/medias/${this.$route.params.id}`
    );

    this.media = res.data.data;
    this.streamReady = this.media.status === "Ready" && Object.keys(this.media.edges).length > 0;
    this.loading = false;

    if (this.streamReady) {
      this.stream();
    }
  }
});
</script>
