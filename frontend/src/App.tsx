import { Calculator } from "./components/Calculator/Calculator";
import { HttpCalculatorApi } from "./api/calculator";

const api = new HttpCalculatorApi();

function App() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900">
      <Calculator api={api} />
    </div>
  );
}

export default App;
