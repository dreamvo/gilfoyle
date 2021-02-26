<template>
  <div>
    <div v-if="mediaId" class="mb-5">
      <v-btn text color="primary" :to="{ name: 'MediaView', id: mediaId }">
        <v-icon size="16" class="mr-1">mdi-keyboard-backspace</v-icon>
        Back
      </v-btn>
    </div>

    <v-card outlined :loading="loading">
      <v-card-title>{{ title }}</v-card-title>

      <v-divider></v-divider>

      <v-col md="12">
        <validation-observer ref="observer" v-slot="{ handleSubmit }">
          <form class="" @submit.prevent="handleSubmit(submit)">
            <validation-provider
                v-slot="{ errors }"
                rules="required|min:1|max:255"
            >
              <v-text-field
                  v-model="form.title"
                  :error-messages="errors"
                  label="Title of the media"
              ></v-text-field>
            </validation-provider>

            <v-btn
                color="primary"
                class="mt-5"
                depressed
                :loading="loading"
                @click="handleSubmit(submit)"
            >
              <v-icon size="16" class="mr-1">mdi-check</v-icon>
              Save
            </v-btn>
          </form>
        </validation-observer>
      </v-col>
    </v-card>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import {setInteractionMode, ValidationObserver, ValidationProvider} from "vee-validate";
import {AxiosResponse} from "axios";
import {DataResponse, Media} from "../types";
import axios from "../services/axios";

setInteractionMode("eager");

interface Data {
  loading: boolean;
  form: {
    title: string;
  };
}

export default Vue.extend({
  name: "MediaForm",
  components: {
    ValidationProvider,
    ValidationObserver
  },
  props: {
    title: {
      type: String,
      default: "Create a new media"
    },
    mediaId: {
      type: String,
      default: ""
    }
  },
  data: (): Data => ({
    loading: true,
    form: {
      title: ""
    }
  }),
  methods: {
    async submit(): Promise<void> {
      this.loading = true

      if (this.mediaId) {
        const res: AxiosResponse<DataResponse<Media>> = await axios.patch(
            `/medias/${this.mediaId}`,
            {
              title: this.form.title
            }
        );

        await this.$router.push({
          name: "MediaView",
          params: {id: res.data.data.id}
        });
        return;
      }

      const res: AxiosResponse<DataResponse<Media>> = await axios.post(
          "/medias",
          {
            title: this.form.title
          }
      );

      await this.$router.push({
        name: "MediaView",
        params: {id: res.data.data.id}
      });
    }
  },
  async created() {
    if (this.mediaId) {
      const res: AxiosResponse<DataResponse<Media>> = await axios.get(
          `/medias/${this.mediaId}`
      );

      this.form = {
        title: res.data.data.title
      };
    }

    this.loading = false;
  }
});
</script>
