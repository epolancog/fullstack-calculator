import { Calculator } from "./components/Calculator/Calculator";
import { HttpCalculatorApi } from "./api/calculator";

const api = new HttpCalculatorApi();

function App() {
  return (
    <div className="bg-gradient-animated min-h-screen flex items-center justify-center p-4">
      <Calculator api={api} />
    </div>
  );
}

export default App;
