<template>
  <v-row>
    <v-col lg="6" md="12" sm="12">
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
        </v-card-subtitle>

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

          <DeleteModal btn-small btn-text="Delete" :action="deleteMedia" />
        </v-card-actions>
      </v-card>
    </v-col>

    <v-col lg="6" md="12" sm="12">
      <v-card flat outlined :loading="loading" :disabled="loading">
        <v-card-title>Streaming</v-card-title>

        <v-card-text
          ><span v-if="!isStreamReady"
            >This media is not yet available for streaming.</span
          >
          <div v-else>
            <VideoPlayer :sources="sources" />
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col cols="12" md="12" v-if="media.edges && media.edges.probe">
      <v-card flat outlined :loading="loading" :disabled="loading">
        <v-card-title>Original file</v-card-title>

        <v-simple-table>
          <template v-slot:default>
            <thead>
              <tr>
                <th class="text-left">
                  Filename
                </th>
                <th class="text-left">
                  Checksum
                </th>
                <th class="text-left">
                  Aspect ratio
                </th>
                <th class="text-left">
                  Resolution
                </th>
                <th class="text-left">
                  Duration
                </th>
                <th class="text-left">
                  Framerate
                </th>
                <th class="text-left">
                  Format
                </th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>{{ media.edges.probe.filename }}</td>
                <td>{{ media.edges.probe.checksum_sha256 }}</td>
                <td>{{ media.edges.probe.aspect_ratio }}</td>
                <td>
                  {{ media.edges.probe.width }}x{{ media.edges.probe.height }}
                </td>
                <td>
                  {{ Math.round(media.edges.probe.duration_seconds) }} sec
                </td>
                <td>{{ Math.round(media.edges.probe.framerate) }} FPS</td>
                <td>{{ media.edges.probe.format }}</td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>
      </v-card>
    </v-col>

    <v-col cols="12" md="12" v-if="media.edges && media.edges.media_files">
      <v-card flat outlined :loading="loading" :disabled="loading">
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
                  <v-chip :color="statusLevel(item.status)" small
                    >{{ item.status }}
                  </v-chip>
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

    <v-col cols="12" md="12">
      <v-card flat outlined :loading="loading" :disabled="loading">
        <v-card-title>Events</v-card-title>

        <v-simple-table>
          <template v-slot:default>
            <thead>
              <tr>
                <th class="text-left">
                  Type
                </th>
                <th class="text-left">
                  Reason
                </th>
                <th class="text-left">
                  Message
                </th>
                <th class="text-left">
                  Date
                </th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>
                  <v-chip small>Normal</v-chip>
                </td>
                <td>MediaCreated</td>
                <td>Media was created through the API</td>
                <td>March 1, 2021 at 6:50 pm</td>
              </tr>
              <tr>
                <td>
                  <v-chip small>Normal</v-chip>
                </td>
                <td>MediaHasUpload</td>
                <td>Original media file was uploaded through the API</td>
                <td>March 1, 2021 at 6:51 pm</td>
              </tr>
              <tr>
                <td>
                  <v-chip small>Normal</v-chip>
                </td>
                <td>EncodingScheduled</td>
                <td>Encoding jobs has been created</td>
                <td>March 1, 2021 at 6:51 pm</td>
              </tr>
              <tr>
                <td>
                  <v-chip small>Error</v-chip>
                </td>
                <td>EncodingFailed</td>
                <td>Encoding job failed for rendition 480p with error : ...</td>
                <td>March 1, 2021 at 6:51 pm</td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { Vue } from "vue-property-decorator";
import { AxiosResponse } from "axios";
import { DataResponse, Media, Source } from "../types";
import axios from "../services/axios";
import DeleteModal from "../components/DeleteModal.vue";
import MediaUploadForm from "../components/MediaUploadForm.vue";
import VideoPlayer from "../components/VideoPlayer.vue";

interface Data {
  poller: number;
  loading: boolean;
  uploadDialog: boolean;
  media: Media;
}

export default Vue.extend({
  components: { VideoPlayer, MediaUploadForm, DeleteModal },
  data: (): Data => ({
    poller: 0,
    loading: true,
    uploadDialog: false,
    media: {} as Media
  }),
  methods: {
    async deleteMedia() {
      await axios.delete(`/medias/${this.$route.params.id}`);
      await this.$router.push({ name: "MediaAll" });
    },
    async fetchMedia() {
      const res: AxiosResponse<DataResponse<Media>> = await axios.get(
        `/medias/${this.$route.params.id}`
      );

      this.media = res.data.data;

      if (!this.isProcessing) {
        clearInterval(this.poller);
      }
    },
    statusLevel(status: string): string {
      const statusValues: { [key: string]: string } = {
        Ready: "success",
        Errored: "error",
        AwaitingUpload: "warning"
      };

      return statusValues[status] || "primary";
    }
  },
  computed: {
    sources(): Source[] {
      return [
        {
          src: `${axios.defaults.baseURL}/medias/${this.media.id}/stream/master.m3u8`,
          type: "application/x-mpegURL"
        }
      ];
    },
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
      return this.isReady || this.media.playable;
    }
  },
  async created() {
    await this.fetchMedia();

    if (this.isProcessing) {
      this.poller = setInterval(this.fetchMedia, 3000);
    }

    this.loading = false;
  }
});
</script>
