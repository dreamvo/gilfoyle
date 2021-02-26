<template>
  <v-row>
    <v-col md="6" sm="12">
      <v-card flat outlined :loading="loading || processing">
        <v-card-title
          >{{ media.title }}
          <v-chip :color="statusLevel" small class="ml-2">{{
            media.status
          }}</v-chip>
        </v-card-title>
        <v-card-subtitle>Created on {{ media.created_at }}</v-card-subtitle>

        <v-card-subtitle class="mt-1">
          <v-text-field
            label="Identifier"
            outlined
            dense
            readonly
            :value="media.id"
          ></v-text-field>
        </v-card-subtitle>

        <v-card-actions>
          <v-dialog v-model="uploadDialog" width="500">
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                v-if="media.status === 'AwaitingUpload'"
                text
                color="primary"
                small
                v-bind="attrs"
                v-on="on"
              >
                <v-icon size="16" dark class="mr-1">
                  mdi-upload
                </v-icon>
                Upload media file
              </v-btn>
            </template>

            <v-card>
              <MediaUploadForm
                title="Upload a media file"
                :media-id="media.id"
              />
            </v-card>
          </v-dialog>

          <v-btn
            small
            text
            color="secondary"
            dark
            class="ml-2"
            :to="{ name: 'MediaUpdate', id: media.id }"
          >
            <v-icon size="16" dark class="mr-1">
              mdi-pencil
            </v-icon>
            Update
          </v-btn>

          <DeleteModal btn-small btn-text="Delete" :action="deleteMedia" />
        </v-card-actions>
      </v-card>

      <v-card
        flat
        outlined
        :loading="loading"
        class="mt-6"
        v-if="media.edges && media.edges.media_files"
      >
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
                <td>{{ item.framerate }} FPS</td>
                <td>{{ Math.round(item.duration_seconds) }} sec</td>
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
          <video ref="player" controls loop width="100%"></video>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { Vue } from "vue-property-decorator";
import { AxiosResponse } from "axios";
import { DataResponse, Media } from "../types";
import axios from "../services/axios";
import Hls from "hls.js";
import DeleteModal from "../components/DeleteModal.vue";
import MediaUploadForm from "../components/MediaUploadForm.vue";

interface Data {
  loading: boolean;
  streamReady: boolean;
  uploadDialog: boolean;
  media: Media;
}

export default Vue.extend({
  components: { MediaUploadForm, DeleteModal },
  data: () => ({
    loading: true,
    streamReady: false,
    uploadDialog: false,
    media: {} as Media
  }),
  methods: {
    stream() {
      const hls = new Hls();
      const stream = `${axios.defaults.baseURL}/medias/${this.media.id}/stream/master.m3u8`;
      const video = this.$refs.player as HTMLMediaElement;

      hls.loadSource(stream);
      hls.attachMedia(video);

      // hls.on(Hls.Events.MANIFEST_PARSED, () => {});
    },
    async deleteMedia() {
      await axios.delete(`/medias/${this.$route.params.id}`);

      await this.$router.push({ name: "MediaAll" });
    }
  },
  computed: {
    statusLevel(): string {
      const status: { [key: string]: string } = {
        Ready: "success",
        Errored: "error",
        AwaitingUpload: "warning"
      };

      return status[this.media.status] || "";
    },
    processing(): boolean {
      return this.media.status === "Processing";
    },
    awaitingUpload(): boolean {
      return this.media.status === "AwaitingUpload";
    }
  },
  async created() {
    const res: AxiosResponse<DataResponse<Media>> = await axios.get(
      `/medias/${this.$route.params.id}`
    );

    this.media = res.data.data;
    this.streamReady =
      this.media.status === "Ready" && Object.keys(this.media.edges).length > 0;
    this.loading = false;

    if (this.streamReady) {
      this.stream();
    }
  }
});
</script>
