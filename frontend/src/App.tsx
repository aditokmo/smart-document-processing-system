import { ReactQueryProvider } from './providers'
import { TanStackRouterProvider } from './providers/TanstackRouterProvider'
import { Toaster } from 'react-hot-toast'

function App() {
  return (
    <ReactQueryProvider>
      <TanStackRouterProvider />
      <Toaster />
    </ReactQueryProvider>
  )
}

export default App
