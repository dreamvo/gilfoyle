<template>
  <v-row>
    <v-col md="6" sm="12">
      <v-card
          flat
          outlined
          :loading="loading || isProcessing"
          :disabled="loading"
      >
        <v-card-title
        >{{ media.title }}
          <v-chip :color="statusLevel(media.status)" small class="ml-2"
          >{{ media.status }}
          </v-chip>
        </v-card-title>
        <v-card-subtitle
        >Created on {{ media.created_at | readableDateHour }}
        </v-card-subtitle
        >

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
                  v-if="awaitingUpload"
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

          <DeleteModal btn-small btn-text="Delete" :action="deleteMedia"/>
        </v-card-actions>
      </v-card>

      <v-card
          flat
          outlined
          :loading="loading"
          :disabled="loading"
          class="mt-6"
          v-if="media.edges && media.edges.media_files"
      >
        <v-card-title>Renditions</v-card-title>

        <v-simple-table>
          <template v-slot:default>
            <thead>
            <tr>
              <th class="text-left">
                Status
              </th>
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
              <td>
                <v-chip :color="statusLevel(item.status)" small>{{ item.status }}</v-chip>
              </td>
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
      <v-card flat outlined :loading="loading" :disabled="loading">
        <v-card-title>Streaming</v-card-title>

        <v-card-text
        ><span v-if="!isStreamReady"
        >This media is not yet available for streaming.</span
        >
          <div v-else class="">
            <VideoPlayer :sources="sources"/>
          </div>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import {Vue} from "vue-property-decorator";
import {AxiosResponse} from "axios";
import {DataResponse, Media, MediaFile, Source} from "../types";
import axios from "../services/axios";
import DeleteModal from "../components/DeleteModal.vue";
import MediaUploadForm from "../components/MediaUploadForm.vue";
import VideoPlayer from "../components/VideoPlayer.vue";

interface Data {
  loading: boolean;
  uploadDialog: boolean;
  media: Media;
  sources: Source[];
}

export default Vue.extend({
  components: {VideoPlayer, MediaUploadForm, DeleteModal},
  data: (): Data => ({
    loading: true,
    uploadDialog: false,
    media: {} as Media,
    sources: []
  }),
  methods: {
    async deleteMedia() {
      await axios.delete(`/medias/${this.$route.params.id}`);

      await this.$router.push({name: "MediaAll"});
    },
    rendition(rendition: MediaFile) {
      this.sources = [
        {
          src: `${axios.defaults.baseURL}/medias/${this.media.id}/stream/${rendition.rendition_name}/${rendition.entry_file}`,
          type: rendition.mimetype,
        }
      ];
    },
    resetRenditions() {
      this.sources = [
        {
          src: `${axios.defaults.baseURL}/medias/${this.media.id}/stream/master.m3u8`,
          type: "application/x-mpegURL"
        }
      ];
    },
    async fetchMedia() {
      const res: AxiosResponse<DataResponse<Media>> = await axios.get(
          `/medias/${this.$route.params.id}`
      );

      this.media = res.data.data;

      if (!this.isProcessing) {
        clearInterval()
      }
    },
    statusLevel(status: string): string {
      const statusValues: { [key: string]: string } = {
        Ready: "success",
        Errored: "error",
        AwaitingUpload: "warning"
      };

      return statusValues[status] || "default";
    },
  },
  computed: {
    isReady(): boolean {
      return this.media.status === "Ready";
    },
    isProcessing(): boolean {
      return this.media.status === "Processing";
    },
    awaitingUpload(): boolean {
      return this.media.status === "AwaitingUpload";
    },
    isStreamReady(): boolean {
      if (this.media.edges.media_files && this.media.edges.media_files.length == 0) {
        return false
      }

      return this.media.status === 'Ready'
    }
  },
  async created() {
    await this.fetchMedia()


    this.resetRenditions();

    if (this.isProcessing) {
      setInterval(this.fetchMedia, 3000)
    }

    this.loading = false;
  }
});
</script>
