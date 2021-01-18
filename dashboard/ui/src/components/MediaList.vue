<template>
  <v-card flat v-if="medias.length" :loading="loading">
    <v-card-title>{{ title }}</v-card-title>
    <v-card-subtitle>{{ medias.length }} items</v-card-subtitle>

    <v-row>
      <v-col md="4" v-for="media in medias" :key="media.id">
        <v-card :to="`/medias/${media.id}`">
          <v-img
            class="white--text align-end"
            height="200px"
            :src="require('@/assets/default_media.jpeg')"
          >
            <v-card-title>{{ media.title }}</v-card-title>
          </v-img>

          <v-card-subtitle class="pb-0">
            Status : {{ media.status }}
          </v-card-subtitle>

          <v-card-subtitle class="pb-0">
            Added on : {{ media.created_at }}
          </v-card-subtitle>

          <v-card-text class="text--primary mt-3"></v-card-text>

          <v-card-actions></v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!--    <v-simple-table>-->
    <!--      <template v-slot:default>-->
    <!--        <thead>-->
    <!--        <tr>-->
    <!--          <th class="text-left">-->
    <!--            Title-->
    <!--          </th>-->
    <!--          <th class="text-left">-->
    <!--            Status-->
    <!--          </th>-->
    <!--          <th class="text-left">-->
    <!--            Creation date-->
    <!--          </th>-->
    <!--        </tr>-->
    <!--        </thead>-->
    <!--        <tbody>-->
    <!--        <tr v-for="media in medias" :key="media.id">-->
    <!--          <td>-->
    <!--            <RouterLink :to="`/medias/${media.id}`"-->
    <!--            >{{ media.title }}-->
    <!--            </RouterLink>-->
    <!--          </td>-->
    <!--          <td>{{ media.status }}</td>-->
    <!--          <td>{{ new Date(media.created_at) }}</td>-->
    <!--        </tr>-->
    <!--        </tbody>-->
    <!--      </template>-->
    <!--    </v-simple-table>-->

    <v-spacer></v-spacer>

    <v-card-actions>
      <v-btn depressed color="#66f" dark @click="loadMore">Load more </v-btn>
    </v-card-actions>
  </v-card>

  <v-card v-else :loading="loading">
    <v-card-title>{{ title }}</v-card-title>

    <v-card-subtitle>There's nothing to show here</v-card-subtitle>
  </v-card>
</template>

<script lang="ts">
import Vue from "vue";
import axios from "../services/axios";
import { AxiosResponse } from "axios";
import { ArrayResponse, Media } from "../types";

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
      const res: AxiosResponse<ArrayResponse> = await axios.get(
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
