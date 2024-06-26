import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import { useEffect, useState } from "react";
import { pingMUP } from "../services/PingService";
import toast from "react-hot-toast";
import DrivingBan from "../models/Mup/DrivingBans";
import DrivingBanList from "../components/Mup/DrivingBan/DrivingBanList";
import { getDrivingBans, getDrivingPermit, getDrivingPermitRequests, getRegistrationRequests, getVehicles, handleRequestPermit } from "../services/MupService";
import Modal from "../components/Shared/Modal/Modal";
import VehicleList from "../components/Mup/Vehicle/VehicleList";
import DrivingPermitCard from "../components/Mup/DrivingPermit/DrivingPermitList";
import DrivingPermit from "../models/Shared/DrivingPermit";
import DrivingPermitsRequestCard from "../components/Mup/DrivingPermitRequest/DrivingPermitRequestList";
import RegistrationRequestCard from "../components/Mup/RegistrationRequest/RegistrationRequestList";
import Registration from "../models/Shared/Registration";
import VehicleDTO from "../models/Mup/VehicleDetails";

const MupPage = () => {
  const navigate = useNavigate();
  const [modalVisible, setModalVisible] = useState(false);
  const [modalContent, setModalContent] = useState<any>(null);
  const [modalHeading, setModalHeading] = useState("");
  
  useEffect(() => {
    const token = localStorage.getItem("token") || "";
    if (token === "") navigate("/");
  });

  const ping = () => {
      pingMUP().then(() => {
        toast.success("Good connection with MUP service");
      }).catch((error) => {
        toast.error("No connection to MUP service");
        console.error(error);
      });
  };

  const checkDrivingBans = () => {
    getDrivingBans().then((drivingBans: DrivingBan[]) => {
      if (drivingBans != null) {
        setModalContent(<DrivingBanList drivingBans={drivingBans} />);
        setModalHeading("Your driving bans");
      } else {
        setModalContent("");
        setModalHeading("No driving bans for your account");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch driving bans");
      console.error(error);
    });
  }

  const checkVehicles = () => {
    getVehicles().then((vehicles: VehicleDTO[]) => {
      if(vehicles != null) {
        setModalContent(<VehicleList vehicles={vehicles} />);
        setModalHeading("Your vehicles");
      }else {
        setModalContent("");
        setModalHeading("No vehicles for your account");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch vehicles");
      console.error(error);
    });
  };

  const checkDrivingPermit = () => { 
    getDrivingPermit().then((drivingPermits: DrivingPermit[]) => {
      if(drivingPermits!= null) {
        setModalContent(<DrivingPermitCard drivingPermits={drivingPermits} />);
        setModalHeading("Your driving permits");
      }else {
        setModalContent(<Button
          label="Request Driving Permit"
          buttonType="button"
          onClick={handleRequestPermit}
      />);
        setModalHeading("No driving permit for your account");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch driving permit");
      console.error(error);
    });
  };

  const checkDrivingPermitRequest = () => {
    getDrivingPermitRequests().then((DrivingPermitRequests: DrivingPermit[]) => {
      if(DrivingPermitRequests!= null) {
        setModalContent(
          <DrivingPermitsRequestCard 
            drivingPermitRequests={DrivingPermitRequests} 
            closeModal={() => setModalVisible(false)} 
            refreshRequests={checkDrivingPermitRequest}
          />
        );
        setModalHeading("List of driving permit requests");
      }else {
        setModalContent("");
        setModalHeading("No pending driving permit requests");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch driving permit requests");
      console.error(error);
    });
  };

  const checkRegistrationRequest = () => {
    getRegistrationRequests().then((RegistrationRequests: Registration[]) => {
      if(RegistrationRequests!= null) {
        setModalContent(
          <RegistrationRequestCard 
            registrationRequests={RegistrationRequests}
            closeModal={() => setModalVisible(false)}
            refreshRequests={checkRegistrationRequest}
          />
        );
        setModalHeading("List of registration requests");
      }else {
        setModalContent("");
        setModalHeading("No pending registration requests");
      }
      setModalVisible(true);
    }).catch(error => {
      toast.error("Failed to fetch registration requests");
      console.error(error);
    });
  };
  
  function closeModal(): void {
    setModalVisible(false);
  };

  return (
    <>
      <HeadingStyled>MUP</HeadingStyled>
      <br />
      <Button buttonType="button" label="Ping Service" onClick={() => ping()} />
      <br />
      <br />
      <Button buttonType="button" label="Check my vehicles" onClick={() => checkVehicles()}/>
      <Button buttonType="button" label="Check my driving bans" onClick={() => checkDrivingBans()}/>
      <Button buttonType="button" label="Check my driving permit" onClick={() => checkDrivingPermit()}/>
      <Button buttonType="button" label="Pending traffic permit requests" onClick={() => checkDrivingPermitRequest()}/>
      <Button buttonType="button" label="Pending registration requests" onClick={() => checkRegistrationRequest()}/>
      <br />
      <br />
      <Button buttonType="button" label="Go Back" onClick={() => window.history.back()}/>
      <br />
      
      <Modal
        heading={modalHeading}
        content={modalContent}
        isVisible={modalVisible}
        onClose={closeModal}
      />
    </>
  );
};

export default MupPage;