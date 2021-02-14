<template>
  <v-card outlined :loading="loading">
    <v-card-title>{{ title }}</v-card-title>
    <v-card-subtitle>{{ medias.length }} items</v-card-subtitle>

    <div v-if="medias.length">
      <v-container>
        <v-row>
          <v-col cols="12" lg="3" md="4" sm="12" v-for="media in medias" :key="media.id">
            <v-card outlined :to="`/medias/${media.id}`">
              <v-img
                  class="white--text align-end"
                  height="200px"
                  :src="require('@/assets/default_media.jpeg')"
              >
                <v-card-title style="display: inline-block;width:95%;white-space: nowrap;overflow: hidden;text-overflow: ellipsis;">
                  {{ media.title }}
                </v-card-title>
              </v-img>

              <v-card-subtitle class="pb-0">
                Status : {{ media.status }}
              </v-card-subtitle>

              <v-card-text class="text--primary mt-3"></v-card-text>

              <v-card-actions></v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>

      <v-spacer></v-spacer>

      <v-card-actions class="justify-center">
        <v-btn class="pl-5 pr-5" depressed color="#66f" dark @click="loadMore"
        >Load more
        </v-btn>
      </v-card-actions>
    </div>
    <v-card-subtitle v-else>There's nothing to show here</v-card-subtitle>
  </v-card>
</template>

<script lang="ts">
import Vue from "vue";
import axios from "../services/axios";
import {AxiosResponse} from "axios";
import {ArrayResponse, Media} from "../types";

interface Data {
  title: string;
  limit: number;
  offset: number;
  medias: Media[];
  loading: boolean;
}

export default Vue.extend({
  name: "MediaList",
  data: (): Data => ({
    title: "Latest medias",
    limit: 6,
    offset: 0,
    medias: [],
    loading: true
  }),
  methods: {
    async loadMore() {
      this.loading = true;
      const medias = await this.fetchMedias();
      this.medias.push(...medias);
      this.loading = false;
      this.offset = this.medias.length;
    },
    async fetchMedias() {
      const res: AxiosResponse<ArrayResponse<Media>> = await axios.get(
          `/medias?limit=${this.limit}&offset=${this.offset}`
      );
      return res.data.data;
    }
  },
  async created() {
    this.medias = await this.fetchMedias();
    this.offset = this.medias.length;
    this.loading = false;
  }
});
</script>
