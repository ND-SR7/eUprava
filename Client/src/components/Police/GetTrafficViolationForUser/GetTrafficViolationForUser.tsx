import React, { useEffect, useState } from "react";
import { TrafficViolationsTable } from "./GetTrafficViolationForUser.styled";
import TrafficViolation from "../../../models/Police/TrafficViolation";
import { getAllTrafficViolationsForUser } from "../../../services/PoliceService";

const GetTrafficViolationForUser: React.FC = () => {
  const [violations, setViolations] = useState<TrafficViolation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchViolations = async () => {
      try {
        const data = await getAllTrafficViolationsForUser();

        if (Array.isArray(data)) {
          setViolations(data);
        } else {
          setViolations([]);
          setError("Unexpected response format.");
          console.error("Unexpected response format:", data);
        }
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchViolations();
  }, []);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error}</p>;

  return (
    <div>
      <h2>Traffic Violations</h2>
      <TrafficViolationsTable>
        <thead>
          <tr>
            <th>JMBG</th>
            <th>Reason</th>
            <th>Description</th>
            <th>Time</th>
            <th>Location</th>
          </tr>
        </thead>
        <tbody>
          {violations.map((violation) => (
            <tr key={violation.id}>
              <td>{violation.violatorJMBG}</td>
              <td>{violation.reason}</td>
              <td>{violation.description}</td>
              <td>{new Date(violation.time).toLocaleString()}</td>
              <td>{violation.location}</td>
            </tr>
          ))}
        </tbody>
      </TrafficViolationsTable>
    </div>
  );
};

export default GetTrafficViolationForUser;
