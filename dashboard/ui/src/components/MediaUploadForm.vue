<template>
  <div>
    <v-card outlined :loading="loading" :disabled="loading">
      <v-card-title>{{ title }}</v-card-title>

      <v-divider></v-divider>

      <v-col md="12">
        <validation-observer ref="observer" v-slot="{ handleSubmit }">
          <form class="" @submit.prevent="handleSubmit(submit)">
            <v-radio-group v-model="form.kind">
              <v-row>
                <v-col cols="12" md="6">
                  <v-card outlined>
                    <v-radio
                        label="Video"
                        value="video"
                        class="ma-3"
                    ></v-radio>
                  </v-card>
                </v-col>


                <v-col cols="12" md="6">
                  <v-card outlined>
                    <v-radio
                        label="Audio only"
                        value="audio"
                        class="ma-3"
                    ></v-radio>
                  </v-card>
                </v-col>
              </v-row>
            </v-radio-group>

            <validation-provider
                v-slot="{ errors, validate }"
                rules="required"
            >
              <v-file-input
                  v-model="form.file"
                  :error-messages="errors"
                  label="Media file"
                  prepend-icon=""
                  show-size
                  outlined
                  @change="validate"
              ></v-file-input>
            </validation-provider>

            <v-card-actions>
              <v-spacer></v-spacer>

              <v-btn
                  color="primary"
                  class="mt-5"
                  depressed
                  :loading="loading"
                  @click="handleSubmit(submit)"
              >
                <v-icon size="16" class="mr-1">mdi-upload</v-icon>
                Upload
              </v-btn>
            </v-card-actions>
          </form>
        </validation-observer>
      </v-col>
    </v-card>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import {setInteractionMode, ValidationObserver, ValidationProvider} from "vee-validate";
import axios from "../services/axios";

setInteractionMode("eager");

interface Data {
  loading: boolean;
  form: {
    kind: string
    file: File | null;
  };
}

export default Vue.extend({
  name: "MediaUploadForm",
  components: {
    ValidationProvider,
    ValidationObserver
  },
  props: {
    title: {
      type: String,
      default: "Upload a media file"
    },
    mediaId: {
      type: String,
      default: ""
    }
  },
  data: (): Data => ({
    loading: true,
    form: {
      kind: 'video',
      file: null,
    }
  }),
  methods: {
    async submit(): Promise<void> {
      this.loading = true

      const formData = new FormData()
      formData.append('file', this.form.file as File)

      await axios.post(
          `/medias/${this.mediaId}/upload/${this.form.kind}`,
          formData,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
      );

      this.$router.go(0)
    }
  },
  async created() {
    this.loading = false;
  }
});
</script>
