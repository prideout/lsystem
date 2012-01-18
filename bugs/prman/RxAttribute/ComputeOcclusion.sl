plugin "GetColorAttr";

class ComputeOcclusion(float samples = 128;)
{
  public void surface(output color Ci, Oi)
  {
    normal Nf = normalize(faceforward(N, I));
    normal Nn = normalize(Nf);
    float occ = occlusion(P, Nn, samples, "maxdist", 100.0);
    //uniform color c = (1,0,1);
    //attribute("user:CrazyColor", c); // This always works.
    uniform color c = GetColorAttr("user:CrazyColor"); // This doesn't work when editing.
    Ci = (1 - occ) * Cs * c * Os;
    Oi = Os;
  }
}
