import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client'
import { setContext } from '@apollo/client/link/context'

// В dev режиме используем proxy из vite.config.js, в prod - абсолютный URL
const graphqlUrl = import.meta.env.PROD 
  ? (import.meta.env.VITE_GRAPHQL_URL || 'http://localhost:8080/graphql')
  : '/graphql'

const httpLink = createHttpLink({
  uri: graphqlUrl,
})

const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem('token')
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '',
    },
  }
})

export const apolloClient = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache(),
})
