<template>
  <div>
    <div v-if="items.length >0">

      <h2>You've already searched for:</h2>
      <b-list-group>
        <b-list-group-item v-for="item in items" :key="item.Url">
          <a :href="`/domains/${item.Url}`">
            <DomainCard  :url="item.Url" :info="item.Info" ></DomainCard>
          </a>

<!--          <b-link :to="`/domains/${item}`" >-->
<!--            http://{{item}}-->
<!--          </b-link>-->
        </b-list-group-item>
      </b-list-group>
    </div>
    <div v-if="items.length ===0">
      <h1>You haven't checked any domains yet!</h1>
      <router-link to="/">Check some in here!</router-link>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import DomainCard from './DomainCard'

export default {
  components: { DomainCard },
  data() {
    return {
      items: []
    }
  },
  methods: {
    getAll() {
      axios({ method: 'GET', url: 'http://localhost:8090/allDomains', headers: { 'content-type': 'text/plain' } })
        .then(result => {
          this.items = result.data.Items
        })
        .catch(error => {
          /*eslint-disable*/
          console.error(error);
          /* eslint-enable */
        })
    }
  },
  mounted () {
    this.getAll()
  }
}
</script>
<style scoped>
  a {
    color: inherit;
  }
  .list-group {
    margin-top: 2rem;
    margin-bottom: 4rem;
  }
</style>
