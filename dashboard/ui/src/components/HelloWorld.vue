<template>
  <v-container>
    <v-row class="">
      <v-col class="mb-4">
        <h1 class="mb-3">
          Health check
        </h1>

        <p>
          utilise un texte en faux latin, le Lorem ipsum ou Lipsum.Wikipedia
          (FR)
        </p>
      </v-col>

      <v-col md="12">
        <h1 class="mb-3">Medias</h1>

        <v-simple-table>
          <template v-slot:default>
            <thead>
              <tr>
                <th class="text-left">
                  Title
                </th>
                <th class="text-left">
                  Status
                </th>
                <th class="text-left">
                  Created at
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="media in medias" :key="media.id">
                <td>
                  <RouterLink :to="`/medias/${media.id}`">{{
                    media.title
                  }}</RouterLink>
                </td>
                <td>{{ media.status }}</td>
                <td>{{ media.created_at }}</td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>

        <v-spacer></v-spacer>
        <div align="center">
          <v-btn depressed color="#66f" dark>Load more</v-btn>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "HelloWorld",
  data: () => ({
    medias: [
      {
        created_at: "2020-12-15T19:56:04.729Z",
        edges: {
          media_files: [
            {
              created_at: "2020-12-15T19:56:04.729Z",
              duration_seconds: 13.25,
              edges: {},
              encoder_preset: "string",
              framerate: 0,
              id: "1139c704-00e9-4c23-8609-e83d0eacb5a5",
              media_type: "video",
              scaled_width: 1280,
              updated_at: "2020-12-15T19:56:04.729Z",
              video_bitrate: 2481
            }
          ]
        },
        id: "1139c704-00e9-4c23-8609-e83d0eacb5a5",
        status: "Processing",
        title: "POLO & PAN — Home Sweet Home (the mixtape)",
        updated_at: "2020-12-15T19:56:04.729Z"
      }
    ]
  })
});
</script>
