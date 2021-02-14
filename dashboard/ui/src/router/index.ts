import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import Home from "../views/Home.vue";
import Medias from "../views/Medias.vue";
import MediaView from "../views/MediaView.vue";
import MediaCreate from "../views/MediaCreate.vue";
import MediaUpdate from "../views/MediaUpdate.vue";
import About from "../views/About.vue";
import NotFound from "../views/NotFound.vue";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "Home",
    component: Home
  },
  {
    path: "/medias",
    name: "MediaAll",
    component: Medias
  },
  {
    path: "/medias/create",
    name: "MediaCreate",
    component: MediaCreate
  },
  {
    path: "/medias/:id/update",
    name: "MediaUpdate",
    component: MediaUpdate
  },
  {
    path: "/medias/:id",
    name: "MediaView",
    component: MediaView
  },
  {
    path: "/about",
    name: "About",
    component: About
  },
  {
    path: "*",
    name: "NotFound",
    component: NotFound
  }
];

const router = new VueRouter({
  mode: "hash",
  base: process.env.BASE_URL,
  routes
});

export default router;
